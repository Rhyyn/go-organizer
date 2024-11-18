import { useState, useEffect } from "react";
// import logo from './assets/images/logo-universal.png';
import logo from "./assets/CLASSES_icons/logo-sram.png";
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

    const handleToggle = () => {
        if (!isActive) {
            ResumeHook();
        } else {
            PauseHook();
        }
        setIsActive(!isActive);
    };

    return (
        <div id="App">
            {/* <img src={logo} id="logo" alt="logo" /> */}
            <div className="menu-container"></div>
            <button className="btn" onClick={getDofusWindows}>
                Refresh
            </button>
            <button onClick={saveOrder}>Save</button>
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
                                        ↑
                                    </button>
                                    <button
                                        className="move-button"
                                        onClick={() => moveDown(index)}
                                        disabled={
                                            index === dofusWindows.length - 1
                                        }
                                    >
                                        ↓
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
            <button onClick={() => logList()}>Log la liste</button>
            <div>
                <label>
                    Organizer is : {isActive ? "Active" : "Paused"}
                    <input
                        className={`custom-checkbox ${
                            isActive ? "active" : "paused"
                        }`}
                        type="checkbox"
                        checked={isActive}
                        onChange={handleToggle}
                    />
                </label>
            </div>
        </div>
    );
}

export default App;
