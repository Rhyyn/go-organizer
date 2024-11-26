import { useState, useEffect, useRef } from "react";
import { TitleBar } from "./TitleBar";
import classIcons from "./ClassIcons";
import upArrow from "./assets/GUI_icons/arrow-up.png";
import downArrow from "./assets/GUI_icons/arrow-down.png";
import reduceWhite from "./assets/GUI_icons/reduceWhite.png";
import expandWhite from "./assets/GUI_icons/expandWhite.png";
import expandRight from "./assets/GUI_icons/expandRight.png";
import expandDown from "./assets/GUI_icons/expandDown.png";
import { EventsOff, EventsOn, WindowSetSize } from "../wailsjs/runtime/runtime";
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
    UpdateTemporaryDofusWindows,
} from "../wailsjs/go/main/App";

function App() {
    const isFirstRun = useRef(true);
    const [isActive, setIsActive] = useState(false);
    const [dofusWindows, setDofusWindows] = useState([]);
    const [keycodes, setKeycodes] = useState([]);

    const [previousKey, setPreviousKey] = useState("");
    const [nextKey, setNextKey] = useState("");
    const [stopOrganizerKey, setStopOrganizerKey] = useState("");
    const [isOnTop, setIsOnTop] = useState(false);

    // First run of the app to get keycodes/keybinds
    useEffect(() => {
        if (isFirstRun.current) {
            getKeyCodes();
            FetchKeybinds();
            isFirstRun.current = false; // Mark as done after the first run
        }
    }, []);

    // TODO: check this
    useEffect(() => {
        EventsOn("updatedCharacterOrder", (newState) => {
            setDofusWindows(newState);
        });

        return () => {
            EventsOff("updatedCharacterOrder");
        };
    }, [isActive]);

    // returns an array of {Code: number, Name: string} from our keycodes.go
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

    // TODO: use event updatedCharacterOrder instead?
    function getDofusWindows() {
        GetDofusWindows().then((result) => {
            if (result !== null) {
                setDofusWindows(result);
                console.log(result);
            }
        });
    }

    // saves order to our characters.ini
    async function saveOrder() {
        await SaveCharacterList(dofusWindows).catch((error) => {
            // need to do this to other functions
            console.error("Failed to save Dofus windows order:", error);
        });
        console.log("updatedorder.. to :");
        console.log(dofusWindows);
    }

    // TODO: Do this auto on fetch
    // ask back to load order of saved characters
    async function loadOrder() {
        console.log("updating order..");
        console.log(dofusWindows);
        await UpdateDofusWindowsOrder(dofusWindows)
            .then((result) => {
                if (result.length != 0) {
                    setDofusWindows(result);
                }
            })
            .catch((error) => {
                console.error("Failed to update Dofus windows order:", error);
            });
        // GetDofusWindows().then(updateWindows);
        // console.log("updatedorder.. to :");
    }

    // TODO: move should change the order
    const moveUp = (index) => {
        if (index > 0) {
            const newList = [...dofusWindows];
            const temp = newList[index];
            newList[index] = newList[index - 1];
            newList[index - 1] = temp;
            setDofusWindows(newList);
            UpdateTemporaryDofusWindows(newList);
        }
    };

    const moveDown = (index) => {
        if (index < dofusWindows.length - 1) {
            const newList = [...dofusWindows];
            const temp = newList[index];
            newList[index] = newList[index + 1];
            newList[index + 1] = temp;
            setDofusWindows(newList);
            UpdateTemporaryDofusWindows(newList);
        }
    };

    // EmitsOn("KeybindsUpdate", () => {
    //     FetchKeybinds();
    // })

    // Fetch keybinds
    const FetchKeybinds = () => {
        GetAllKeyBindings().then((result) => {
            console.log(result);
            Object.values(result).map((keybind) => {
                switch (keybind.Action) {
                    case "StopOrganizer":
                        setStopOrganizerKey(keybind.KeyName);
                        break;
                    case "NextChar":
                        setNextKey(keybind.KeyName);
                        break;
                    case "PreviousChar":
                        setPreviousKey(keybind.KeyName);
                        break;
                    default:
                        break;
                }
            });
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

    // Observer to know when organizer is active or not
    useEffect(() => {
        EventsOn("updateOrganizerRunningState", (newState) => {
            setIsActive(newState);
        });

        return () => {
            EventsOff("updateOrganizerRunningState");
        };
    }, [isActive]);

    // 1 == full mode, 2 = Overlay mode
    const [isWindowFull, setIsWindowFull] = useState(true);
    const [windowFullPosition, setWindowFullPosition] = useState({});
    const [windowOverlayPosition, setWindowOverlayPosition] = useState({});

    const handleWindowMode = async () => {
        if (isWindowFull) {
            WindowSetSize(200, 46);
            setIsWindowFull(false);
        } else {
            WindowSetSize(392, 800);
            setIsWindowFull(true);
        }
    };

    const handleOverlayListDirection = () => {};

    const [activeChar, setActiveChar] = useState(null);
    EventsOn("CharSelectedEvent", (activeChar) => {
        console.log(activeChar);
        setActiveChar(activeChar);
    });

    const handleAlwaysOntop = () => {
        SetAlwaysOnTop(!isOnTop);
        setIsOnTop(!isOnTop);
    };

    return (
        <div id="App">
            {isWindowFull ? (
                <div className="full-mode">
                    <TitleBar></TitleBar>
                    <div className="menu-container">
                        <button
                            className={`btn ${isOnTop ? "on-top" : ""}`}
                            onClick={handleAlwaysOntop}
                        >
                            Pin to top
                        </button>
                        <button className="btn" onClick={getDofusWindows}>
                            Fetch
                        </button>
                        <button className="btn" onClick={loadOrder}>
                            Load
                        </button>
                        <button className="btn" onClick={saveOrder}>
                            Save
                        </button>
                    </div>
                    <button
                        onClick={() => handleWindowMode()}
                        className="button-to-overlay-mode"
                    >
                        Switch to Overlay Mode
                    </button>
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
                                                onClick={() =>
                                                    WinActivate(window.hwnd)
                                                }
                                            >
                                                {window.CharacterName}
                                            </div>
                                            {/* |
                                            <div className="class-name">
                                                {window.Class}
                                            </div> */}
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
                                                    index ===
                                                    dofusWindows.length - 1
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
                                            parseInt(
                                                selectedOption.dataset.key
                                            ),
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
                                value={previousKey.toLowerCase()}
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
                                            parseInt(
                                                selectedOption.dataset.key
                                            ),
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
                                        value={key.Name.toLowerCase()}
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
                                            parseInt(
                                                selectedOption.dataset.key
                                            ),
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
            ) : (
                <div className="overlay-mode" style={{ widows: "1" }}>
                    <img
                        src={expandWhite}
                        className="expand-icon"
                        alt="expand window"
                        title="Full Mode"
                        onClick={() => handleWindowMode()}
                        style={{ widows: "2" }}
                    ></img>
                    {dofusWindows.length === 0 ? (
                        <span>
                            No Dofus windows found. Use Fetch in full mode
                        </span>
                    ) : (
                        <div className="overlay-characters-container">
                            {dofusWindows.map((window, index) => (
                                <div
                                    // fix this, it's not working properly
                                    className={`overlay-character-item ${
                                        activeChar === window.hwnd
                                            ? "char-active"
                                            : "char-inactive"
                                    }`}
                                    key={window.CharacterName}
                                    style={{ widows: "2" }}
                                >
                                    <img
                                        className="overlay-class-icon"
                                        src={
                                            classIcons[
                                                window.Class.toLowerCase()
                                            ] || classIcons.default
                                        }
                                        alt={`${
                                            window.Class || "Default"
                                        } Icon`}
                                        style={{ widows: "2" }}
                                        onClick={() => WinActivate(window.hwnd)}
                                    ></img>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}

export default App;
