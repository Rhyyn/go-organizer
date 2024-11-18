export namespace main {
	
	export class WindowInfo {
	    title: string;
	    hwnd: number;
	    CharacterName: string;
	    Class: string;
	    Order: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.hwnd = source["hwnd"];
	        this.CharacterName = source["CharacterName"];
	        this.Class = source["Class"];
	        this.Order = source["Order"];
	    }
	}

}

