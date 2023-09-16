let unqId = parseInt(Math.random() * 500_000) // The key sent to the server to get the specific simulation back
let width = parseInt(window.innerWidth/2) // The width of the simulation

const virLength = 1000; // How large the area is in the simulation
const realLength = virLength/width // The area mapped to the client's screen
let simTime = 0; // Counts the years
let simData = []; // This is a list of frames to show on the screen
let sheepSize = 25;
let wolfSize = 25;
let grassSize = 25;
let readoutTimer = -1;
const year = 120; // How long a year is
let started = false;
let requesting = false;
//let serverAvailable = true;

async function makeRequest(subDomain, requestData) {
    //while (!serverAvailable) {

    //}
    //serverAvailable = false;

    if (requestData != "") {
        requestData = "/" + requestData;
    }
    const response = await fetch("http://localhost:9090/" + subDomain + "/" + unqId + "" + requestData);
    let data = await response.json()

    //serverAvailable = true;
    return data;
}

function findFrame(frames) {
    let lowest = simTime * 2;
    let index = 0;
    for (let i = 0; i < frames.length; i++) {
        if (frames[i][0] < lowest) {
            lowest = frames[i][0]
            index = i;
        }
    }
    return frames[index];
}

async function setup() {
    for (let i = 1; i < document.body.childNodes.length; i++) {
        try {
        let child = document.body.childNodes[i];
        child.style.opacity = "0";
        } catch {}
    }
    let children = document.querySelectorAll(".startInfo");
    for (let i = 0; i < children.length; i++) {
        let child = children[i];
        child.style.opacity = "1";
    }

    const cnv = createCanvas(width, width); // The little window in the middle of the simulation
    simData.push(await makeRequest("start", "")); // Add the frame to the queue
    frameRate(25);// Keep conistency, rather low and consistent
}

// Makes a request to the server, gets a large amount to keep the buffer full
async function getData() {
    if (requesting) return;
    requesting = true;
    simData = simData.concat(await makeRequest("sim", ""));
    requesting = false;
}

function checkInRange(value, min, max) {
    return value >= min && value <= max
}

async function draw() {
    if (!started) return;

    // The readout timer
    if (readoutTimer != -1) {
        if (readoutTimer == 0) {
            buttonReadout.onclick();
        }
        readoutTimer--;
    }

    // Get the data
    data = findFrame(simData); // Get a frame
    simData.shift(); // Remove the first index of the queue
    simTime++;
    
    // Get the lists of each object
    const grass = data[1];
    const wolf = data[2];
    const sheep = data[3];

    // Show information
    divStats.innerHTML = "Grass: " + grass.length + "<br>Wolf: " + wolf.length + "<br>Sheep: " + sheep.length + "<br>Total: " + (grass.length + wolf.length + sheep.length) + "<br>Time Spent: " + Math.round(simTime / year) + " Years";

    background(128, 64, 0); // Makes the simulation screen brown

    fill(0, 255, 0); // Green
    for (let i = 0; i < grass.length; i++) {
        square(grass[i][0]/realLength, grass[i][1]/realLength, width*grassSize/virLength);
    }

    fill(128, 128, 128); // Grey
    for (let i = 0; i < wolf.length; i++) {
        square(wolf[i][0]/realLength, wolf[i][1]/realLength, width*wolfSize/virLength);
    }

    fill(240, 240, 240); // White
    for (let i = 0; i < sheep.length; i++) {
        square(sheep[i][0]/realLength, sheep[i][1]/realLength, width*sheepSize/virLength);
    }

    // Keep a second of buffer to allow smooth framerate
    if (simData.length < 24) {
        await getData();
        return;
    }
}

btnStart.onclick = async () => {
    started = true;
    
    for (let i = 1; i < document.body.childNodes.length; i++) {
        try {
        let child = document.body.childNodes[i];
        child.style.opacity = "1";
        } catch {}
    }
    let children = document.querySelectorAll(".startInfo");
    for (let i = 0; i < children.length; i++) {
        let child = children[i];
        child.remove();
    }
}

 // Scrolls to and from the help page, the scrollbar needs to be shown to scroll, but can be instntly removed
 buttonHelp.onclick = () => {
    document.body.style.overflow = "visible";
    window.scrollTo(0, window.innerHeight);
    document.body.style.overflow = "hidden";
}

buttonBack.onclick = () => {
    document.body.style.overflow = "visible";
    window.scrollTo(0, 0);
    document.body.style.overflow = "hidden";
}

buttonRestart.onclick = async () => {
    // Delete simulation
    let data = await makeRequest("remove", "")

    // Get a new ID
    unqId = parseInt(inputId.value);
    if (unqId != unqId) {
        unqId = parseInt(Math.random() * 500_000);
    }

    data = await makeRequest("exists", "")
    if (data == "\"true\"") {
        data = await makeRequest("sim", "")
    }
    else {
        // Start a new simulation
        data = await makeRequest("start", "")
    }

    simData = [data]; // get rid of queue
    simTime = 0; // Start the timer again
}

buttonReadout.onclick = async () => {
    // Get the data
    const data = await makeRequest("readout", "");
    divReadout.innerHTML = data + "<br>Time: " + (simTime/year);

    // Check if the client has asked for a timed readout
    let val = timerTime.value;
    if (val == "") return;
    try {
        val = parseInt(val);
    }
    catch {
        return;
    }
    // Keep it within 10 minutes and make sure it's positive
    if (val > 600 || val < 0) return;
    readoutTimer = val * 25;
}

buttonConfig.onclick = async () => {
    // Get config data from the DOM
    const sheepData = divControlsLeft.children[0].children;
    const wolfData = divControlsRight.children[1].children;
    const grassData = divControlsRight.children[2].children;

    buttonConfig.style.color = "red";

    let newConfig = []; // Create an array to send
    let temp;

    // Sheep
    temp = sheepData[2].children[0].value; // Get the number from the input
    if (!checkInRange(temp, 1, 100)) return; // Make sure it's valid
    newConfig.push(parseFloat(temp)); // Add it to the array

    temp = sheepData[4].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    sheepSize = parseInt(temp)
    newConfig.push(sheepSize);

    temp = sheepData[6].children[0].value;
    if (!checkInRange(temp, 0, 1)) return;
    newConfig.push(parseFloat(temp));

    temp = sheepData[8].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[10].children[0].value;
    if (!checkInRange(temp, 1, 100000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[12].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[14].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[16].children[0].value;
    if (!checkInRange(temp, 1, 1000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[18].children[0].value;
    if (!checkInRange(temp, 1, 1000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[20].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[22].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[24].children[0].value;
    if (!checkInRange(temp, 1, 150)) return;
    newConfig.push(parseFloat(temp));

    // Wolves
    temp = wolfData[2].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    newConfig.push(parseFloat(temp));

    temp = wolfData[4].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    wolfSize = parseInt(temp);
    newConfig.push(wolfSize);

    temp = wolfData[6].children[0].value;
    if (!checkInRange(temp, 0, 1)) return;
    newConfig.push(parseFloat(temp));

    temp = wolfData[8].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[10].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[12].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[14].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[16].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[18].children[0].value;
    if (!checkInRange(temp, 1, 1000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[20].children[0].value;
    if (!checkInRange(temp, 1, 150)) return;
    newConfig.push(parseFloat(temp));

    // Grass
    temp = grassData[2].children[0].value;
    if (!checkInRange(temp, 1, 20000)) return;
    newConfig.push(parseFloat(temp));

    temp = grassData[4].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    grassSize = parseInt(temp)
    newConfig.push(grassSize);

    temp = grassData[6].children[0].value;
    if (!checkInRange(temp, 1, 200)) return;
    newConfig.push(parseFloat(temp));

    temp = grassData[8].children[0].value;
    if (!checkInRange(temp, 1, 2000)) return;
    newConfig.push(parseFloat(temp));

    temp = grassData[10].children[0].value;
    if (!checkInRange(temp, 1, 1000)) return;
    newConfig.push(parseFloat(temp));

    // Send it to the server
    console.log(await makeRequest("config", JSON.stringify(newConfig))); // Should get "Success"
    buttonConfig.style.color = "white";
}