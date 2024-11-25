import React from "react";
import logo from "./assets/images/appIcon.png";
import crossIcon from "./assets/GUI_icons/cross-white.png";
import { Quit } from "../wailsjs/runtime/runtime";

export const TitleBar = () => {
    return (
        <header>
            <div className="title-bar-container" style={{ widows: "1" }}>
                <img
                    src={logo}
                    className="app-logo"
                    alt="logo"
                    data-wails-drag
                />
                <div className="app-title-container" data-wails-drag>
                    Go-organizer
                </div>
                <img
                    onClick={() => Quit()}
                    className="cross-icon"
                    src={crossIcon}
                    alt="quit"
                />
            </div>
        </header>
    );
};
