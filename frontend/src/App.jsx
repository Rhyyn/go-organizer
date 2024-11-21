import { useState, useEffect } from "react";
// import logo from './assets/images/logo-universal.png';
import classIcons from "./ClassIcons";
import upArrow from "./assets/GUI_icons/arrow-up.png";
import downArrow from "./assets/GUI_icons/arrow-down.png";
import { EventsOff, EventsOn } from "../wailsjs/runtime/runtime";
import "./App.css";
import {
    GetDofusWindows,
    UpdateDofusWindowsOrder,
    PauseHook,
    ResumeHook,
    GetKeycodes,
    SaveKeybind,
    GetAllKeyBindings,
    SaveCharacterList,
    WinActivate,
    SetAlwaysOnTop,
    StopHook,
} from "../wailsjs/go/main/App";

function App() {
    const [isFirst, setIsFirst] = useState(true);
    const [isActive, setIsActive] = useState(false);
    const [dofusWindows, setDofusWindows] = useState([]);
    const [keycodes, setKeycodes] = useState([]);

    const [previousKey, setPreviousKey] = useState("");
    const [nextKey, setNextKey] = useState("");
    const [stopOrganizerKey, setStopOrganizerKey] = useState("");

    // EventsOn("KeybindsUpdate", (updatedKeybinds) => {
    //     console.log("event received", updatedKeybinds);
    //     setStopOrganizerKey(updatedKeybinds["StopOrganizer"].KeyName);
    //     setNextKey(updatedKeybinds["NextChar"].KeyName);
    //     setPreviousKey(updatedKeybinds["PreviousChar"].KeyName);
    // });

    useEffect(() => {
        EventsOn("updatedCharacterOrder", (newState) => {
            setDofusWindows(newState);
        });

        return () => {
            EventsOff("updatedCharacterOrder");
        };
    }, [isActive]);

    function getKeyCodes() {
        GetKeycodes().then((result) => {
            const keycodesArray = Object.entries(result).map(
                ([code, name]) => ({
                    Code: parseInt(code, 10),
                    Name: name,
                })
            );
            const sortedList = keycodesArray.sort((a, b) => a.Code - b.Code);
            setKeycodes(sortedList);
        });
    }

    function getDofusWindows() {
        GetDofusWindows().then((result) => {
            if (result !== null) {
                setDofusWindows(result);
                console.log(result);
            }
        });
    }

    async function saveOrder() {
        console.log("updating order..");
        await SaveCharacterList(dofusWindows).catch((error) => {
            // If Go returned an error, handle it here
            console.error("Failed to save Dofus windows order:", error);
        });
        // GetDofusWindows().then(updateWindows);
        console.log("updatedorder.. to :");
        console.log(dofusWindows);
    }

    async function loadOrder() {
        console.log("updating order..");
        await UpdateDofusWindowsOrder(dofusWindows)
            .then((result) => {
                setDofusWindows(result);
            })
            .catch((error) => {
                // If Go returned an error, handle it here
                console.error("Failed to update Dofus windows order:", error);
            });
        // GetDofusWindows().then(updateWindows);
        console.log("updatedorder.. to :");
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

    useEffect(() => {
        if (isFirst) {
            getKeyCodes();
            FetchKeybinds();
        }
        setIsFirst(false);
    }, [isFirst]);

    // Fetch keybinds
    const FetchKeybinds = () => {
        GetAllKeyBindings().then((result) => {
            setStopOrganizerKey(result["StopOrganizer"].KeyName.toUpperCase());
            setNextKey(result["NextChar"].KeyName.toUpperCase());
            setPreviousKey(result["PreviousChar"].KeyName.toUpperCase());
        });
    };

    // Call backend to save new keybind then fetch
    async function saveKeybinds(keycode, keyname, keybindName) {
        await SaveKeybind(keycode, keyname, keybindName).then(() => {
            FetchKeybinds();
        });
    }

    const handleActiveToggle = () => {
        if (!isActive) {
            ResumeHook();
        } else {
            PauseHook();
        }
        setIsActive(!isActive);
    };

    useEffect(() => {
        EventsOn("updateMainHookState", (newState) => {
            setIsActive(newState);
        });

        return () => {
            EventsOff("updateMainHookState");
        };
    }, [isActive]);

    return (
        <div id="App">
            {/* <img src={logo} id="logo" alt="logo" /> */}
            <div className="menu-container">
                <button className="btn" onClick={getDofusWindows}>
                    Fetch
                </button>
                <button className="btn" onClick={SetAlwaysOnTop}>
                    Pin to top
                </button>
                <button className="btn" onClick={loadOrder}>
                    Order
                </button>
                <button className="btn" onClick={saveOrder}>
                    Save
                </button>
                <button className="btn" onClick={StopHook}>
                    StopHook
                </button>
            </div>
            <div id="dofusWindowList">
                {dofusWindows.length === 0 ? (
                    <p>No Dofus windows found. Use Fetch</p>
                ) : (
                    <div className="windows-container">
                        {dofusWindows.map((window, index) => (
                            <div className="window-item" key={index}>
                                <div className="left-container">
                                    <div
                                        className="character-name"
                                        title={`Click to activate ${window.CharacterName}`}
                                        onClick={() => WinActivate(window.hwnd)}
                                    >
                                        {window.CharacterName}
                                    </div>
                                    |
                                    <div className="class-name">
                                        {window.Class}
                                    </div>
                                </div>
                                <div className="icon-container">
                                    <img
                                        className="class-icon"
                                        src={
                                            classIcons[
                                                window.Class.toLowerCase()
                                            ] || classIcons.default
                                        }
                                        alt={`${
                                            window.Class || "Default"
                                        } Icon`}
                                    ></img>
                                </div>
                                <div className="move-buttons-container">
                                    <button
                                        title="Move Up"
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
                                        title="Move Down"
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
            <div className="toggle-container">
                <div className="up-toggle-container">
                    <label
                        className={`toggle-label ${
                            isActive ? "active" : "paused"
                        }`}
                    >
                        Organizer is : {isActive ? "Active" : "Off"}
                        <input
                            className={`custom-checkbox ${
                                isActive ? "active" : "paused"
                            }`}
                            type="checkbox"
                            checked={isActive}
                            onChange={handleActiveToggle}
                        />
                    </label>
                </div>
            </div>
            <div className="bottom-container">
                <label className="dropdown-label">
                    Toggle Organizer :
                    <select
                        className="dropdown"
                        value={stopOrganizerKey}
                        onKeyDown={(event) => {
                            // might need some more
                            event.preventDefault();
                        }}
                        onChange={async (event) => {
                            const selectedOption =
                                event.target.options[
                                    event.target.selectedIndex
                                ];

                            try {
                                await saveKeybinds(
                                    parseInt(selectedOption.dataset.key),
                                    selectedOption.value.toLowerCase(),
                                    "StopOrganizer"
                                );
                            } catch (error) {
                                console.error(
                                    "Error saving keybind or fetching keys:",
                                    error
                                );
                            }
                        }}
                    >
                        {keycodes.map((key) => (
                            <option
                                key={key.Code}
                                value={key.Name}
                                data-key={key.Code}
                                className="dropdown-option"
                            >
                                {key.Name}
                            </option>
                        ))}
                    </select>
                </label>
                <label className="dropdown-label">
                    Previous Character :
                    <select
                        className="dropdown"
                        value={previousKey}
                        onKeyDown={(event) => {
                            event.preventDefault();
                        }}
                        onChange={async (event) => {
                            const selectedOption =
                                event.target.options[
                                    event.target.selectedIndex
                                ];

                            try {
                                await saveKeybinds(
                                    parseInt(selectedOption.dataset.key),
                                    selectedOption.value.toLowerCase(),
                                    "PreviousChar"
                                );
                            } catch (error) {
                                console.error(
                                    "Error saving keybind or fetching keys:",
                                    error
                                );
                            }
                        }}
                    >
                        {keycodes.map((key) => (
                            <option
                                key={key.Code}
                                value={key.Name}
                                data-key={key.Code}
                                className="dropdown-option"
                            >
                                {key.Name}
                            </option>
                        ))}
                    </select>
                </label>
                <label className="dropdown-label">
                    Next Character :
                    <select
                        className="dropdown"
                        value={nextKey}
                        onKeyDown={(event) => {
                            event.preventDefault();
                        }}
                        onChange={async (event) => {
                            const selectedOption =
                                event.target.options[
                                    event.target.selectedIndex
                                ];

                            try {
                                await saveKeybinds(
                                    parseInt(selectedOption.dataset.key),
                                    selectedOption.value.toLowerCase(),
                                    "NextChar"
                                );
                            } catch (error) {
                                console.error(
                                    "Error saving keybind or fetching keys:",
                                    error
                                );
                            }
                        }}
                    >
                        {keycodes.map((key) => (
                            <option
                                key={key.Code}
                                value={key.Name}
                                data-key={key.Code}
                                className="dropdown-option"
                            >
                                {key.Name}
                            </option>
                        ))}
                    </select>
                </label>
            </div>
        </div>
    );
}

export default App;
