<script lang="ts" setup>
import {ref, reactive, watch } from 'vue'
import {OpenSerialPort, SerialClose, GetPorts, SetPortSetting, GetReceiveData, SendData, SetDelimiter} from '../../wailsjs/go/main/App'
import {main}  from '../../wailsjs/go/models'

enum Direction {
  SEND = 'S',   // Send
  RECEIVE = 'R' // Receive
}

// シリアルポートの受信結果を取得するインターバル
let serialReadInterval: number;

// LogRowというインターフェースを定義
interface LogRow {
    direction: Direction;
    commBytes: number[];
}
const commLog = ref<Array<LogRow>>([]);

// 送信データ
const dataString = ref("");
// 表示モード
const viewMode = ref("HEX");

// 受信データ区切り方法
const delimitType = ref("byTimer");
const delimitInterval = ref("150");
const delimitCode = ref("0");

// COMポートの状態(オープン・クローズ)
const comOpenStatus = ref(false);

// COMポート設定
const portSetting = reactive<main.PortSetting>({
  portNo: 1,
  baud: 9600,
  parity: "N",
  stopBit: 1,
  wordLength: 8,
});
SetPortSetting(portSetting);

// 区切り方法
const delimitSetting = {
  delimitByCode: false,
  code: 0x00,
  intervalMs: 150
};
SetDelimiter(delimitSetting);

// COMポート
const portNoList = ref<Array<number>>([]);

// ボーレート
const baudList = ref<Array<number>>(
  [2400, 4800, 9600, 19200, 38400, 57600, 115200]
);

// パリティー
const parityList = reactive([
  {name:"NONE", value:"N"},
  {name:"偶数", value:"E"},
  {name:"奇数", value:"O"}
]);

// ストップビット長
const stopBitList = ref<Array<number>>([1,2]);

// ワード長
const wordLengthList = ref<Array<number>>([7, 8]);

  /**
   * 「接続」をクリック
   */
const onClickOpenSerial = () => {
  comOpenStatus.value = true;

  OpenSerialPort().then(
    result => {
      if (!result) {  // オープン失敗
        comOpenStatus.value = false;
        alert("COMオープンエラー");
        return;
      }

      // オープン成功
      comOpenStatus.value = true;

      // 受信インターバル
      serialReadInterval = setInterval(() => {

        GetReceiveData().then(  // 受信した
          result => {
            if (result == null) return;

            for (let i=0; i<result.length; i++) {
              let targetRow: LogRow = {direction: Direction.RECEIVE, commBytes:[]};
              
              // 追加の文字列を取得
              if (result[i].append) {  // 追加受信
                // 最後に受信したログ行を取得
                for (let j=commLog.value.length-1; j>=0; j--) {
                  if (commLog.value[j].direction == Direction.RECEIVE) {
                    targetRow = commLog.value[j];
                    targetRow.commBytes = targetRow.commBytes.concat(result[i].data);
                    break;
                  }
                }

                // バイト数がゼロの場合は追加対象が見つからず新規と判断する
                if (targetRow.commBytes.length <= 0) {
                  targetRow.commBytes = targetRow.commBytes.concat(result[i].data);
                  commLog.value.push(targetRow);
                }
              }
              else {  // 新規受信
                targetRow = {direction: Direction.RECEIVE, commBytes: result[i].data};
                commLog.value.push(targetRow);
              }
            }

            // ログの一番下までスクロール
            scrollBottomLog();
          }
        );
      }, 500);
    }
  );
}

/**
 * 「送信」をクリック
 */
const onClickSendData = () => {
  SendData(dataString.value).then(
    result => {
      if (result.errorMessage != '') {
        alert(result.errorMessage);
        return;
      }

      // 送信結果をログへ追加
      let targetRow: LogRow = {
        direction: Direction.SEND,
        commBytes: result.data
      }
      commLog.value.push(targetRow);

      // ログの一番下までスクロール
      scrollBottomLog();
    }
  )
}

/**
 * 「切断」をクリック
 */
const onClickCloseSerial = () => {
  // データの読み込み処理を停止
  clearInterval(serialReadInterval);

  SerialClose().then(
    () => {
      comOpenStatus.value = false;
    }
  )
}

// ポート一覧を取得
GetPorts().then(
  result => {
    for (let i=0; i<result.portCount; i++) {
      portNoList.value.push(result.portNumbers[i]);
    }
    portSetting.portNo = result.portNumbers[0];
  }
);

// ポートの設定を変更
watch(portSetting, (newPs, oldPs) => {
  SetPortSetting(newPs);
});

/**
 * ログ行をテキスト形式へ整形する
 * @param logrow 
 */
const formatLogRowAsText = (logrow: LogRow) => {
  let buf = "";

  for (let i=0; i < logrow.commBytes.length; i++) {
    switch (viewMode.value) {
      case "HEX":
        buf += (buf == "" ? '' : ' ') +
              toHex(logrow.commBytes[i]);
        break;
      case "DEC":
        buf += (buf == "" ? '' : ' ') +
              logrow.commBytes[i].toString(10);
        break;
      case "STR":
        buf += numberToAscii(logrow.commBytes[i]);
        break;
    }
  }
  return buf;
}

/**
 * 区切り文字の「変更」クリック
 */
const onClickSetDelimit = () => {
  if (delimitType.value == "byTimer") {
    if (Number.isInteger(delimitInterval.value)) {
      alert("インターバルの指定に誤りがあります。");
      return;
    }

    delimitSetting.delimitByCode = false;
    delimitSetting.intervalMs = Number(delimitInterval.value);
    SetDelimiter(delimitSetting);
  }
  else {
    const codeValue = delimitCode.value.startsWith('0x') ? parseInt(delimitCode.value, 16) : parseInt(delimitCode.value, 10);

    // codeValue が数値ではない、または 1 から 255 の範囲外の場合に警告
    if (!Number.isInteger(codeValue) || codeValue <= 0 || codeValue > 255) {
      alert("区切りコードの指定に誤りがあります。");
      return;
    }

    delimitSetting.delimitByCode = true;
    delimitSetting.code = Number(delimitCode.value);
    SetDelimiter(delimitSetting);
  }
}

/**
 * 「ログクリア」をクリック
 */
const onLogClear = ()  => {
  commLog.value = [];
}

/**
 * 該当するASCIIコードへ変換する。
 * 文字として表現できるコードでなければ代りに'・'を返す。
 * @param value 
 */
function numberToAscii(value: number): string {
  // ASCIIコードが表示可能な文字かどうかをチェック
  if (value >= 32 && value <= 126) {
      return String.fromCharCode(value);
  } else {
      return '・'; // 表示不可能な場合は「・」を返す
  }
}  

/**
 * 数値を1バイトのHEX文字列で返す
 * @param value 
 */
function toHex(value: number): string {
  const hex = value.toString(16);
  return hex.padStart(2, '0').toUpperCase();
}

/**
 * ログ表示エリアを最終行までスクロール
 */
function scrollBottomLog() {
  const commDataElement = document.getElementById('comm_data');
  if (commDataElement != null) {
    commDataElement.scrollTop = commDataElement.scrollHeight;
  }
};
</script>

<template>
  <div id="container">
    <div id="control-panel">
      <div id="connection">
        接続：
        <button type="button" @click="onClickOpenSerial()" v-bind:disabled="comOpenStatus">接続</button>
        <button type="button" @click="onClickCloseSerial()" v-bind:disabled="!comOpenStatus">切断</button>
      </div>
      <hr>
      ポート：
      <select v-model="portSetting.portNo"  v-bind:disabled="comOpenStatus">
        <option v-for="port in portNoList" v-bind:value="port">
          COM{{ port }}
        </option>
      </select>
      <hr>
      ボーレート：
      <select v-model="portSetting.baud">
        <option v-for="baud in baudList" v-bind:value="baud">
          {{ baud }}bps
        </option>
      </select>
      <hr>
      パリティ：
      <select v-model="portSetting.parity">
        <option v-for="parity in parityList" v-bind:value="parity.value">
          {{  parity.name }}
        </option>
      </select>
      <hr>
      ストップビット：
      <label v-for="stopBit in stopBitList">
        <input type="radio" name="stopbit" v-bind:value="stopBit" v-model="portSetting.stopBit">{{ stopBit }}&nbsp;
      </label>
      <hr>
      ワード長：
      <label v-for="wordLength in wordLengthList">
        <input type="radio" name="wordLength" v-bind:value="wordLength" v-model="portSetting.wordLength">{{ wordLength }}&nbsp;
      </label>
    </div>
    <div id="comlog_panel">
      <div id="send_data">
        <input type="text" v-model="dataString">
        <button @click="onClickSendData()" v-bind:disabled="!comOpenStatus">送信</button>
      </div>
      <div id="control_code">
        <button @click="dataString += ' 0x1B'">ESC</button>
        <button @click="dataString += ' 0x0D'">CR</button>
        <button @click="dataString += ' 0x0A'">LF</button>
        <button @click="dataString += ' 0x02'">STX</button>
        <button @click="dataString += ' 0x03'">ETX</button>
        <button @click="dataString += ' 0x06'">ACK</button>
        <button @click="dataString += ' 0x15'">NAK</button>
        <button @click="dataString += ' 0x1A'">EOF</button>
      </div>
      <div id="input-explanation">
        コードの指定は数値の先頭に、16進数は"0x"、8進数"0"、2進数は”0b"を付けます。10進数は数値の先頭に何も付けません。<br>
        文字列は半角シングルクォーテーション（"）で囲んでください。<br>
        コードは半角スペースで区切りってください。<br>
        例： 0x02&nbsp;0x31&nbsp;012&nbsp;0b1111&nbsp;"ABC TEST"&nbsp;0x0d&nbsp;0x0a
      </div>
      <div id="select_view_mode">
        表示形式：
        <label><input type="radio" name="viewmode" value="HEX" v-model="viewMode">16進数</label>
        <label><input type="radio" nam e="viewmode" value="DEC" v-model="viewMode">10進数</label>        
        <label><input type="radio" name="viewmode" value="STR" v-model="viewMode">文字列</label>
      </div>
      <div>
        区切り方法：
        <label>
          <input type="radio" name="delimit" value="byTimer" v-model="delimitType">タイマー
          <input type="text" name="delimit_interval" v-model="delimitInterval" style="width:4em;text-align:right;" v-select-all>ms
        </label>
        &nbsp;&nbsp;&nbsp;
        <label>
          <input type="radio" name="delimit" value="byCode" v-model="delimitType">コード
          <input type="text" name="delimit_code" v-model="delimitCode" style="width:4em;text-align:right;" v-select-all>
        </label>
        &nbsp;&nbsp;&nbsp;
        <button @click="onClickSetDelimit()">変更</button>
      </div>
      <div>
        <button @click="onLogClear()" style="width:8em;color:red;">ログクリア</button>
      </div>
      <div id="comm_data">
        <ul>
          <li v-for="logrow in commLog" v-bind:class="[{send : logrow.direction == Direction.SEND},{receive : logrow.direction == Direction.RECEIVE}, {logrow}, ]">
            {{ formatLogRowAsText(logrow) }}
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style>
div {
  box-sizing: border-box;
  text-align: left;
}

ul,li {
  margin:0;
  padding:0;
}

/* 画面全体のコンテナー設定 */
#container {
  display: flex;
  flex-wrap: nowrap; /* フレックスアイテムを折り返さない */
  padding-top: 1em;
  height: 100%;
  min-width: 0; /* 子要素がオーバーフローするのを防ぐ */
}

/* 左側のコントロールパネル */
#control-panel {
  flex-shrink: 0; /* 幅が縮小されないように設定 */
  padding-left: 1em;
  padding-right: 1em;
  width: 18em; /* 固定幅 */
  min-width: 18em; /* 最小幅も設定 */
}

/* 接続・切断のボタン */
#connection button {
  margin-right: 1em;
}

/* 送受信データの表示エリア */
#comlog_panel {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  overflow: auto; /* コンテンツが多い場合にスクロール可能に */
}

/* データの送信項目 */
#send_data {
  border: none;
  border-bottom: none;
  line-height: 2.5em;
  height: 2.5em;
  width: 100%;
  padding-left: .5em;
  padding-right: .5em;
  display: flex;
  align-items: center;
}

#send_data input {
  flex-grow: 1;
  height: 1.5em;
  font-family: monospace;
}

#send_data button {
  height: 2em;
}

#control_code {
  padding-left: 1em;
  min-height: 2em;
}

#control_code button {
  margin-left: .1em;
  width: 4em;
}

/* 表示モードの選択項目 */
#select_view_mode {
  padding-left: 1em;
  min-height: 2em;
  line-height: 2em;
}

#select_view_mode label {
  margin-right: 1em;
  display: inline-block;
}

/* 送受信データの表示エリア */
#comm_data {
  flex-grow: 1;
  overflow: auto; /* 内容が多い場合にスクロール可能に */
  color: white;
  padding: 0.2em;
  border: solid 1px white;
  white-space: pre-wrap; /* 長い文字列でも適切に折り返し */
  font-family: monospace;
}

/* 送受信行 */
.logrow {
  list-style-type: none;
}

.logrow.receive::before {
  content: "受信:";
}

.logrow.send::before {
  content: "送信:";
}

#input-explanation {
  font-family: monospace;
}

</style>