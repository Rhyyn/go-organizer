import { useState, useEffect } from "react";
// import logo from './assets/images/logo-universal.png';
import logo from "./assets/CLASSES_icons/logo-sram.png";
import upArrow from "./assets/GUI_icons/arrow-up.png";
import downArrow from "./assets/GUI_icons/arrow-down.png";
import { EventsEmit, EventsOff, EventsOn } from "../wailsjs/runtime/runtime";
import "./App.css";
import {
    GetDofusWindows,
    UpdateDofusWindowsOrder,
    PauseHook,
    ResumeHook,
    GetKeycodes,
    SaveStopOrgaKeyBind,
    SavePreviousCharKeybind,
    SaveNextCharKeybind,
    GetAllKeyBindings,
    SaveCharacterList,
} from "../wailsjs/go/main/App";

function App() {
    const [isFirst, setIsFirst] = useState(true);
    const [isActive, setIsActive] = useState(false);
    const [isShiftActive, setIsShiftActive] = useState(false);
    const [isAltActive, setIsAltActive] = useState(false);
    const [dofusWindows, setDofusWindows] = useState([]);
    const [keycodes, setKeycodes] = useState([]);

    const [previousKey, setPreviousKey] = useState("");
    const [nextKey, setNextKey] = useState("");
    const [stopOrganizerKey, setstopOrganizerKey] = useState("");

    function fetchSavedKeys() {
        GetAllKeyBindings().then((result) => {
            console.log(result["StopOrganizer"].KeyName);
            setstopOrganizerKey(result["StopOrganizer"].KeyName);
            setNextKey(result["NextChar"].KeyName);
            setPreviousKey(result["PreviousChar"].KeyName);
        });
    }

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

    // const updateWindows = (windows) => {
    //     setDofusWindows(windows); // Update the state with the list of windows
    // };

    useEffect(() => {
        if (isFirst) {
            getKeyCodes();
            fetchSavedKeys();
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
                    Refresh
                </button>
                <button className="btn" onClick={loadOrder}>
                    Load
                </button>
                <button className="btn" onClick={saveOrder}>
                    Save
                </button>
                {/* <button className="btn">Previous</button>
                <button className="btn">Next</button> */}
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
                    {/* <button
                        className={`btn2 ${
                            isToggleListening
                                ? "toggle-active"
                                : "toggle-paused"
                        }`}
                        onClick={handleToggleKeybind}
                    >
                        {defaultKeybindText}
                    </button> */}
                </div>
            </div>
            <div className="bottom-container">
                <label className="dropdown-label">
                    Stop Organizer :
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
                                await SaveStopOrgaKeyBind(
                                    parseInt(selectedOption.dataset.key),
                                    selectedOption.value
                                );

                                fetchSavedKeys();
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
                                value={key.Name} // Set the value to key.Name
                                data-key={key.Code} // Set data-key to key.Code
                            >
                                {key.Name}
                            </option>
                        ))}
                    </select>
                </label>
                <label className="dropdown-label">
                    Previous :
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
                                await SavePreviousCharKeybind(
                                    parseInt(selectedOption.dataset.key),
                                    selectedOption.value
                                );

                                fetchSavedKeys();
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
                                value={key.Name} // Set the value to key.Name
                                data-key={key.Code} // Set data-key to key.Code
                            >
                                {key.Name}
                            </option>
                        ))}
                    </select>
                </label>
                <label className="dropdown-label">
                    Next :
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
                                await SaveNextCharKeybind(
                                    parseInt(selectedOption.dataset.key),
                                    selectedOption.value
                                );

                                fetchSavedKeys();
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
                                value={key.Name} // Set the value to key.Name
                                data-key={key.Code} // Set data-key to key.Code
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
