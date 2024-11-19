import { useState, useEffect } from "react";
// import logo from './assets/images/logo-universal.png';
import logo from "./assets/CLASSES_icons/logo-sram.png";
import upArrow from "./assets/GUI_icons/arrow-up.png";
import downArrow from "./assets/GUI_icons/arrow-down.png";
import { EventsEmit, EventsOff, EventsOn } from "../wailsjs/runtime/runtime";
import "./App.css";
import {
    GetDofusWindows,
    UpdateDofusWindows,
    UpdateDofusWindowsOrder,
    PauseHook,
    ResumeHook,
} from "../wailsjs/go/main/App";

function App() {
    const [isFirst, setIsFirst] = useState(true);
    const [isActive, setIsActive] = useState(false);
    const [isShiftActive, setIsShiftActive] = useState(false);
    const [isAltActive, setIsAltActive] = useState(false);

    function getDofusWindows() {
        UpdateDofusWindows().then((result) => {
            const sortedList = [...result].sort((a, b) => a.Order - b.Order);
            setDofusWindows(sortedList);
        });
    }

    function saveOrder() {
        UpdateDofusWindowsOrder(dofusWindows);
        GetDofusWindows().then(updateWindows);
        console.log("updating order..");
        console.log(dofusWindows);
    }

    const moveUp = (index) => {
        if (index > 0) {
            const newList = [...dofusWindows];
            const temp = newList[index];
            newList[index] = newList[index - 1];
            newList[index - 1] = temp;
            setDofusWindows(newList);
            console.log(newList);
        }
    };

    const moveDown = (index) => {
        if (index < dofusWindows.length - 1) {
            const newList = [...dofusWindows];
            const temp = newList[index];
            newList[index] = newList[index + 1];
            newList[index + 1] = temp;
            setDofusWindows(newList);
            console.log(newList);
        }
    };

    const [dofusWindows, setDofusWindows] = useState([]);

    const updateWindows = (windows) => {
        setDofusWindows(windows); // Update the state with the list of windows
    };

    useEffect(() => {
        if (isFirst) {
            GetDofusWindows()
                .then(updateWindows)
                .catch((error) => {
                    console.error("Error fetching Dofus windows:", error); // Error handling
                });
        }
        setIsFirst(false);
    }, [isFirst]);

    const logList = () => {
        console.log(dofusWindows);
    };

    const handleActiveToggle = () => {
        if (!isActive) {
            ResumeHook();
        } else {
            PauseHook();
        }
        setIsActive(!isActive);
    };

    const handleModifierToggle = (modifier) => {
        if (modifier === "shift") {
            setIsShiftActive(!isShiftActive);
        } else if (modifier === "alt") {
            setIsAltActive(!isAltActive);
        }
    };

    const [defaultKeybindText, setDefaultKeybindText] = useState("Set Keybind");

    useEffect(() => {
        EventsOn("statusUpdated", (newKeybind) => {
            setDefaultKeybindText(newKeybind);
        });

        return () => {
            EventsOff("statusUpdated");
        };
    }, []);

    const [isToggleListening, setIsToggleListening] = useState(false);
    const handleToggleKeybind = () => {
        if (isActive) {
            PauseHook();
            setIsActive(!isActive);
        }
        setDefaultKeybindText("Listening for input..");
        setIsToggleListening(!isToggleListening);

        if (!isActive) {
            ResumeHook();
            setIsActive(!isActive);
        }
    };

    return (
        <div id="App">
            {/* <img src={logo} id="logo" alt="logo" /> */}
            <div className="menu-container">
                <button className="btn" onClick={getDofusWindows}>
                    Refresh
                </button>
                <button className="btn" onClick={saveOrder}>
                    Save
                </button>
                <button className="btn">Previous</button>
                <button className="btn">Next</button>
            </div>
            {/* <div className="sub-menu-container">
                <label className="sub-btn-text">
                    Shift
                    <input
                        type="checkbox"
                        checked={isShiftActive}
                        onChange={() => handleModifierToggle("shift")}
                    />
                </label>
                <label className="sub-btn-text">
                    Alt
                    <input
                        type="checkbox"
                        checked={isAltActive}
                        onChange={() => handleModifierToggle("alt")}
                    />
                </label>
            </div> */}
            <div id="dofusWindowList">
                {dofusWindows.length === 0 ? (
                    <p>No Dofus windows found.</p>
                ) : (
                    <div className="windows-container">
                        {dofusWindows.map((window, index) => (
                            <div className="window-item" key={index}>
                                <div className="left-container">
                                    <div className="character-name">
                                        {window.CharacterName}
                                    </div>
                                    |
                                    <div className="class-name">
                                        {window.Class}
                                    </div>
                                </div>
                                <div className="move-buttons-container">
                                    <button
                                        className="move-button"
                                        onClick={() => moveUp(index)}
                                        disabled={index === 0}
                                    >
                                        <img
                                            className="move-button-arrow"
                                            src={upArrow}
                                            id="downArrow"
                                            alt="downArrow"
                                        />
                                    </button>
                                    <button
                                        className="move-button"
                                        onClick={() => moveDown(index)}
                                        disabled={
                                            index === dofusWindows.length - 1
                                        }
                                    >
                                        <img
                                            className="move-button-arrow"
                                            src={downArrow}
                                            id="upArrow"
                                            alt="upArrow"
                                        />
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
            <button onClick={() => logList()}>Log la liste</button>
            <div className="toggle-container">
                <div className="up-toggle-container">
                    <label
                        className={`toggle-label ${
                            isActive ? "active" : "paused"
                        }`}
                    >
                        Organizer is : {isActive ? "Active" : "Paused"}
                        <input
                            className={`custom-checkbox ${
                                isActive ? "active" : "paused"
                            }`}
                            type="checkbox"
                            checked={isActive}
                            onChange={handleActiveToggle}
                        />
                    </label>
                    <button
                        className={`btn2 ${
                            isToggleListening
                                ? "toggle-active"
                                : "toggle-paused"
                        }`}
                        onClick={handleToggleKeybind}
                    >
                        {defaultKeybindText}
                    </button>
                </div>
            </div>
        </div>
    );
}

export default App;
