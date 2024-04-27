export namespace main {
	
	export class ComPorts {
	    portNumbers: number[];
	    portCount: number;
	
	    static createFrom(source: any = {}) {
	        return new ComPorts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.portNumbers = source["portNumbers"];
	        this.portCount = source["portCount"];
	    }
	}
	export class DelimiterSetting {
	    intervalMs: number;
	    delimitByCode: boolean;
	    code: number;
	
	    static createFrom(source: any = {}) {
	        return new DelimiterSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.intervalMs = source["intervalMs"];
	        this.delimitByCode = source["delimitByCode"];
	        this.code = source["code"];
	    }
	}
	export class PortSetting {
	    portNo: number;
	    baud: number;
	    parity: string;
	    stopBit: number;
	    wordLength: number;
	
	    static createFrom(source: any = {}) {
	        return new PortSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.portNo = source["portNo"];
	        this.baud = source["baud"];
	        this.parity = source["parity"];
	        this.stopBit = source["stopBit"];
	        this.wordLength = source["wordLength"];
	    }
	}
	export class ReceiveData {
	    append: boolean;
	    data: number[];
	
	    static createFrom(source: any = {}) {
	        return new ReceiveData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.append = source["append"];
	        this.data = source["data"];
	    }
	}
	export class SendData {
	    data: number[];
	    errorMessage: string;
	
	    static createFrom(source: any = {}) {
	        return new SendData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.errorMessage = source["errorMessage"];
	    }
	}

}

