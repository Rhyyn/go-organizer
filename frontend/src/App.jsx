import { useState, useEffect } from "react";
// import logo from './assets/images/logo-universal.png';
import logo from "./assets/CLASSES_icons/logo-sram.png";
import "./App.css";
import { SetWindowForeground } from "../wailsjs/go/main/App";
import { GetDofusWindows } from "../wailsjs/go/main/App";
import { UpdateDofusWindows } from "../wailsjs/go/main/App";

function App() {
    const [isFirst, setIsFirst] = useState(true);
    // const [resultText, setResultText] = useState(
    //     "Please enter your name below 👇"
    // );
    // const [name, setName] = useState("");
    // const updateName = (e) => setName(e.target.value);
    // const updateResultText = (result) => setResultText(result);

    // function greet() {
    //     Greet(name).then(updateResultText);
    // }

    // function call() {
    //     TestReturn().then((result) => {
    //         setResultText(result);
    //     });
    // }

    function getDofusWindows() {
        UpdateDofusWindows().then((result) => {
            setDofusWindows(result);
        });
    }

    const moveUp = (index) => {
        if (index > 0) {
            const newList = [...dofusWindows];
            const temp = newList[index];
            newList[index] = newList[index - 1];
            newList[index - 1] = temp;
            setDofusWindows(newList);
        }
    };

    const moveDown = (index) => {
        if (index < dofusWindows.length - 1) {
            const newList = [...dofusWindows];
            const temp = newList[index];
            newList[index] = newList[index + 1];
            newList[index + 1] = temp;
            setDofusWindows(newList);
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

    return (
        <div id="App">
            {/* <img src={logo} id="logo" alt="logo" /> */}
            <div className="menu-container"></div>
            <button className="btn" onClick={getDofusWindows}>
                Refresh
            </button>
            <button>Save</button>
            <button>Load</button>
            <button>Delete</button>
            <button>Précédent</button>
            <button>Suivant</button>
            <div id="dofusWindowList">
                {dofusWindows.length === 0 ? (
                    <p>No Dofus windows found.</p>
                ) : (
                    <ul>
                        {dofusWindows.map((window, index) => (
                            <li key={index}>
                                <strong>{window.title}</strong> - Hwnd:{" "}
                                {window.hwnd}
                                <div>
                                    {/* Arrows to move up and down */}
                                    <button
                                        onClick={() => moveUp(index)}
                                        disabled={index === 0}
                                    >
                                        ↑
                                    </button>
                                    <button
                                        onClick={() => moveDown(index)}
                                        disabled={
                                            index === dofusWindows.length - 1
                                        }
                                    >
                                        ↓
                                    </button>
                                </div>
                            </li>
                        ))}
                    </ul>
                )}
            </div>
            <button onClick={() => SetWindowForeground(dofusWindows[0].hwnd)}>
                Foreground
            </button>
        </div>
    );
}

export default App;
