let unqId = parseInt(Math.random() * 500_000) // The key sent to the server to get the specific simulation back
let width = parseInt(window.innerWidth/2)

const virLength = 1000; // How large the area is in the simulation
const realLength = virLength/width // The area mapped to the client's screen
let simTime = 0;
let simData = [];
let sheepSize = 25;
let wolfSize = 25;
let grassSize = 25;

async function setup() {
    const cnv = createCanvas(width, width); // The little window in the middle of the simulation
    const response = await fetch("http://localhost:9090/start/" + unqId);
    let data = await response.json()
    simData.push(data);
    frameRate(24);// Keep conistency, rather low and consistent
}

// Makes a request to the server, gets a large amount to keep the buffer full
async function getData() {
    const response = await fetch("http://localhost:9090/sim/" + unqId);
    const data = await response.json();
    simData = simData.concat(data);
}

async function draw() {
    // Get the data
    data = simData[0];
    simData.shift();
    simTime++;
    
    const grass = data[0];
    const wolf = data[1];
    const sheep = data[2];

    // Show information
    divStats.innerHTML = "Grass: " + grass.length + "<br>Wolf: " + wolf.length + "<br>Sheep: " + sheep.length + "<br>Total: " + (grass.length + wolf.length + sheep.length) + "<br>Time Spent: " + Math.round(simTime / 3600) + " Years";

    background(128, 64, 0); // Makes the little window blank

    fill(0, 255, 0); // Green
    for (let i = 0; i < grass.length; i++) {
        square(grass[i][0]/realLength, grass[i][1]/realLength, width*grassSize/1000);
    }

    fill(128, 128, 128); // Grey
    for (let i = 0; i < wolf.length; i++) {
        square(wolf[i][0]/realLength, wolf[i][1]/realLength, width*wolfSize/1000);
    }

    fill(240, 240, 240); // White
    for (let i = 0; i < sheep.length; i++) {
        square(sheep[i][0]/realLength, sheep[i][1]/realLength, width*sheepSize/1000);
    }

    // Keep a second of buffer to allow smooth framerate
    if (simData.length < 24) {
        await getData();
        return;
    }
}

let onHelp = 1;

buttonHelp.onclick = () => {
    helpText.style.opacity = onHelp.toString();
    onHelp = 0;
    document.body.style.overflow = "visible";
    window.scrollTo(0, window.innerHeight);
    document.body.style.overflow = "hidden";
}

buttonBack.onclick = () => {
    helpText.style.opacity = onHelp.toString();
    onHelp = 1;
    document.body.style.overflow = "visible";
    window.scrollTo(0, 0);
    document.body.style.overflow = "hidden";
}

buttonRestart.onclick = async () => {
    let response = await fetch("http://localhost:9090/remove/" + unqId);
    let data = await response.json();

    unqId = parseInt(Math.random() * 500_000); // New ID
    response = await fetch("http://localhost:9090/start/" + unqId);
    data = await response.json()

    simData = [data]; // get rid of queue
    simTime = 0; // Start the timer again
}

buttonReadout.onclick = async () => {
    const response = await fetch("http://localhost:9090/readout/" + unqId);
    const data = await response.json();
    divReadout.innerHTML = data;
}

buttonConfig.onclick = async () => {
    const sheepData = divControlsRight.children[1].children;
    const wolfData = divControlsLeft.children[0].children;
    const grassData = divControlsLeft.children[1].children;

    let newConfig = [];
    let temp;

    // Sheep
    temp = sheepData[1].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    newConfig.push(parseFloat(temp));

    temp = sheepData[3].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    sheepSize = parseInt(temp)
    newConfig.push(sheepSize);

    temp = sheepData[5].children[0].value;
    if (!checkInRange(temp, 0, 1)) return;
    newConfig.push(parseFloat(temp));

    temp = sheepData[7].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[9].children[0].value;
    if (!checkInRange(temp, 1, 100000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[11].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[13].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[15].children[0].value;
    if (!checkInRange(temp, 1, 1000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[17].children[0].value;
    if (!checkInRange(temp, 1, 1000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[19].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[21].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = sheepData[23].children[0].value;
    if (!checkInRange(temp, 1, 150)) return;
    newConfig.push(parseFloat(temp));

    // Wolves
    temp = wolfData[1].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    newConfig.push(parseFloat(temp));

    temp = wolfData[3].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    wolfSize = parseInt(temp);
    newConfig.push(wolfSize);

    temp = wolfData[5].children[0].value;
    if (!checkInRange(temp, 0, 1)) return;
    newConfig.push(parseFloat(temp));

    temp = wolfData[7].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[9].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[11].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[13].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[15].children[0].value;
    if (!checkInRange(temp, 1, 10000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[17].children[0].value;
    if (!checkInRange(temp, 1, 1000)) return;
    newConfig.push(parseInt(temp));

    temp = wolfData[19].children[0].value;
    if (!checkInRange(temp, 1, 150)) return;
    newConfig.push(parseFloat(temp));

    // Grass
    temp = grassData[1].children[0].value;
    if (!checkInRange(temp, 1, 20000)) return;
    newConfig.push(parseFloat(temp));

    temp = grassData[3].children[0].value;
    if (!checkInRange(temp, 1, 100)) return;
    grassSize = parseInt(temp)
    newConfig.push(grassSize);

    temp = grassData[5].children[0].value;
    if (!checkInRange(temp, 1, 200)) return;
    newConfig.push(parseFloat(temp));

    temp = grassData[7].children[0].value;
    if (!checkInRange(temp, 1, 2000)) return;
    newConfig.push(parseFloat(temp));

    temp = grassData[9].children[0].value;
    if (!checkInRange(temp, 1, 1000)) return;
    newConfig.push(parseFloat(temp));

    let finalString = JSON.stringify(newConfig);
    const response = await fetch("http://localhost:9090/config/" + unqId + "/" + finalString);
    const data = await response.json();
    console.log(data);
}

function checkInRange(value, min, max) {
    return value >= min && value <= max
}