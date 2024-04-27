package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// パソコンのシリアルポート一覧
type ComPorts struct {
	PortNumbers [20]uint32 `json:"portNumbers"`
	PortCount   uint32     `json:"portCount"`
}

// シリアルポートの設定情報
type PortSetting struct {
	PortNo     uint32 `json:"portNo"`
	Baud       uint32 `json:"baud"`
	Parity     string `json:"parity"`
	StopBit    byte   `json:"stopBit"`
	WordLength byte   `json:"wordLength"`
}

// データ区切り情報
type DelimiterSetting struct {
	IntervalMs    uint32 `json:"intervalMs"`
	DelimitByCode bool   `json:"delimitByCode"`
	Code          uint16 `json:"code"`
}

// 受信データ
type ReceiveData struct {
	Append bool     `json:"append"`
	Data   []uint16 `json:"data"`
}

// 送信データ
type SendData struct {
	Data         []uint16 `json:"data"`
	ErrorMessage string   `json:"errorMessage"`
}

// App struct
type App struct {
	ctx        context.Context
	portHandle syscall.Handle

	// シリアルポートの設定
	ps PortSetting

	// 区切り方法の設定
	ds DelimiterSetting

	// データ受信処理用
	receiveData   chan ReceiveData
	cancelRequest chan struct{}
	cancelDone    chan struct{}
	delimiter     chan DelimiterSetting
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.ds = DelimiterSetting{IntervalMs: 150, DelimitByCode: false, Code: 0x00}
}

// シリアルポートのパラメータを指定する。
func (a *App) SetPortSetting(ps PortSetting) {
	var dcb DCB
	a.ps = ps

	if a.portHandle != NULL { // 既にポートが開いていたら設定を行う
		// 現在のシリアルポートの設定を取得
		GetCommState(a.portHandle, &dcb)
		// ボーレート
		dcb.BaudRate = a.ps.Baud
		// ワード長
		dcb.ByteSize = a.ps.WordLength
		// パリティ
		switch a.ps.Parity {
		case "N": // なし
			dcb.Parity = NOPARITY
		case "O": // 奇数
			dcb.Parity = ODDPARITY
		case "E": // 偶数
			dcb.Parity = EVENPARITY
		default:
			dcb.Parity = NOPARITY
		}
		// ストップビット
		if a.ps.StopBit == 1 {
			dcb.StopBits = ONESTOPBITS
		} else {
			dcb.StopBits = TWOSTOPBITS
		}
		// シリアルポート設定
		SetCommState(a.portHandle, &dcb)
	}
}

func (a *App) SetDelimiter(d DelimiterSetting) {
	a.ds = d
	if a.delimiter == nil { // まだチャンネルが生成されていなければ無視
		return
	}

	// 変更を送る
	a.delimiter <- d
}

// シリアルポートをオープンする。
// 事前に指定されたパラメータに基づいてシリアルポートの設定を行う。
func (a *App) OpenSerialPort() bool {
	var comName = "COM" + strconv.FormatInt(int64(a.ps.PortNo), 10)

	pCharComName, _ := syscall.BytePtrFromString(comName)
	a.portHandle, _ = CreateFile(pCharComName, GENERIC_READ|GENERIC_WRITE, 0, nil, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, 0)

	if a.portHandle == syscall.InvalidHandle { // ポートオープンに失敗
		return false
	}

	// シリアルポートを設定
	a.SetPortSetting(a.ps)

	// シリアルポート受信処理を走らせる
	a.cancelRequest = make(chan struct{})
	a.cancelDone = make(chan struct{})
	a.receiveData = make(chan ReceiveData, 128)
	a.delimiter = make(chan DelimiterSetting, 1)

	// 区切り情報(初期値)を送る
	a.delimiter <- a.ds

	// シリアルポート受信処理
	go func() {
		var ticker *time.Ticker = nil
		var buf []uint16
		var delimiterSetting DelimiterSetting
		var ok bool

		shouldAppend := false

		for {
			var errors uint32
			var commStat COMMSTAT
			var numberOfGot uint32

			// 区切り文字の処理
			select {
			case delimiterSetting, ok = <-a.delimiter: // 区切り情報を変更
				if ok { // チャンネルが開いている
					if delimiterSetting.DelimitByCode { // コードによる区切り
						if ticker != nil {
							ticker.Stop()
						}
						// コードによる区切りの場合は150ms固定
						ticker = time.NewTicker(time.Duration(150) * time.Millisecond)
					} else { // タイマーによる区切り
						if ticker != nil {
							ticker.Stop()
						}
						ticker = time.NewTicker(time.Duration(delimiterSetting.IntervalMs) * time.Millisecond)
					}
				}

			default:
			}

			// 受信処理
			// 受信バイト数をエラー情報と共に取得
			ClearCommError(a.portHandle, &errors, &commStat)
			// シリアルポートからデータを受信
			temp := make([]uint8, commStat.CbInQue)
			if commStat.CbInQue > 0 {
				ReadFile(a.portHandle, &temp[0], commStat.CbInQue, &numberOfGot, nil)
				// byte -> uint16 変換
				for _, b := range temp {
					buf = append(buf, uint16(b))
				}
			}

			select {
			case <-ticker.C: // 定期受信処理を実行
				if delimiterSetting.DelimitByCode { // 受信したコードでデータを区切る
					// データ受信が無ければ以下の処理は行わない
					if len(buf) <= 0 {
						continue
					}

					var t ReceiveData
					t.Append = shouldAppend
					for len(buf) > 0 {
						// 1要素取り出す
						w := buf[0]
						t.Data = append(t.Data, w)
						buf = buf[1:]

						// 区切りコード??
						if w == delimiterSetting.Code {
							// 区切りコードまでのデータをフロントエンドへ渡す
							a.receiveData <- t
							t.Data = []uint16{}
						}
					}

					// 区切りコードがバッファに無かったか、区切りコードの後もデータが続いていた
					if len(t.Data) > 0 {
						a.receiveData <- t
						shouldAppend = true
					} else { // 丁度よく区切りコードで終了
						shouldAppend = false
					}

				} else {
					// データ受信が無ければ以下の処理は行わない
					if len(buf) <= 0 {
						shouldAppend = false
						continue
					}

					// データをフロントエンドへ渡す
					var t ReceiveData
					t.Append = shouldAppend
					t.Data = append(t.Data, buf...)
					a.receiveData <- t

					// バッファをクリア
					buf = []uint16{}

					// 次も時間内に受信があれば継続受信として扱う
					shouldAppend = true
				}

			case <-a.cancelRequest: // 受信処理の終了要求が来た
				if ticker != nil {
					ticker.Stop()
				}
				// 処理を終了したことを通知
				close(a.cancelDone)
				return

			default:
			}
		}
	}()

	return true
}

// バッファに貯めた受信データを取得する
func (a *App) GetReceiveData() []ReceiveData {
	var dlist []ReceiveData

	// シリアルポートからの受信データを取得
	for {
		select {
		case t := <-a.receiveData: // チャンネルのデータを取り出す
			dlist = append(dlist, t)
		default: // チャンネルが空になったら受信データ引き渡し
			return dlist
		}
	}
}

// シリアルポートを閉じる。
func (a *App) SerialClose() {
	// ポートが閉じていれば何もしない
	if a.portHandle == NULL {
		return
	}

	// 区切り文字をセットするチャンネルを破棄
	close(a.delimiter)
	a.delimiter = nil

	// 受信処理の終了要求
	close(a.cancelRequest)

	// 受信処理が終了するまで待つ
	<-a.cancelDone

	// シリアルポートを閉じる
	CloseHandle(a.portHandle)
	a.portHandle = NULL
}

// GetPorts は、利用可能なCOMポートの一覧を取得する。
// GetPorts は、ComPorts構造体を返す。
// この構造体には、利用可能なCOMポートの番号とその総数が含まれる。
func (a *App) GetPorts() ComPorts {
	var portinfo ComPorts

	// COMポートの一覧を取得
	GetCommPorts(&portinfo.PortNumbers[0], uint32(20), &portinfo.PortCount)

	return portinfo
}

// 入力された内容に従ってシリアルポートからデータを送信する
func (a *App) SendData(sendStr string) SendData {
	// var written uint32
	var senddata SendData
	var written uint32

	// 入力内容をバイト配列へ変換
	data, error := a.parseInput(sendStr)
	if len(data) > 0 { // 入力内容にエラーがなければシリアルポートからデータを送信
		WriteFile(a.portHandle, &data[0], uint32(len(data)), &written, nil)
	}

	// 送信した内容をフロントエンドへ返す
	senddata.ErrorMessage = ""
	if error != nil {
		senddata.ErrorMessage = error.Error()
	}
	for _, c := range data {
		senddata.Data = append(senddata.Data, uint16(c))
	}

	return senddata
}

// 入力内容をバイト文字列へ変換する
func (a *App) parseInput(input string) ([]byte, error) {
	var result []byte
	var inString bool
	var currentString strings.Builder

	// 入力を一文字ずつ処理する
	for i := 0; i < len(input); i++ {
		c := input[i]

		if inString {
			if c == '"' {
				// 文字列モードの終了
				inString = false
				str := currentString.String()
				for _, char := range str {
					if char > 127 {
						return nil, fmt.Errorf("ASCII文字列のみ指定可能です: %v", char)
					}
					result = append(result, byte(char))
				}
				currentString.Reset()
			} else {
				// 文字列モードでの文字追加
				currentString.WriteByte(c)
			}
		} else {
			if c == '"' {
				// 文字列モードの開始
				inString = true
			} else if c != ' ' {
				// 数値の解析開始
				start := i
				for i < len(input) && input[i] != ' ' {
					i++
				}
				token := input[start:i]
				i--

				val, err := strconv.ParseInt(token, 0, 10)
				if err != nil {
					return nil, fmt.Errorf("数値の指定に誤りがあります: %v", err)
				}
				if val > 255 || val < 0 {
					return nil, fmt.Errorf("数値の範囲は0～255です: %v", val)
				}
				result = append(result, byte(val))
			}
		}
	}

	// 閉じダブルクォートがない場合の文字列処理
	if currentString.Len() > 0 {
		str := currentString.String()
		for _, char := range str {
			if char > 127 {
				return nil, fmt.Errorf("non-ASCII character detected: %v", char)
			}
			result = append(result, byte(char))
		}
	}

	return result, nil
}
