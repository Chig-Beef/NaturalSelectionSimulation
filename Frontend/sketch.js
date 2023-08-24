let unqId = parseInt(Math.random() * 500_000) // The key sent to the server to get the specific simulation back
let width = parseInt(window.innerWidth/2)

const virLength = 1000;
const realLength = virLength/width
let simTime = 0;

async function setup() {
    const cnv = createCanvas(width, width); // The little window in the middle of the simulation
    const response = await fetch("http://localhost/start/" + unqId);
    const data = await response.json();
    frameRate(24);
}

// Performance testing
let buckets = []
for (let i = 0; i < 100; i++) {
    buckets.push(0);
}

async function draw() {
    let t1 = Date.now();

    // Get the data
    const response = await fetch("http://localhost/sim/" + unqId);
    let data;
    try {data = await response.json();}
    catch {
        return;
    }

    simTime++;
    
    const grass = data[0];
    const wolf = data[1];
    const sheep = data[2];

    // Show information
    divStats.innerHTML = "Grass: " + grass.length + "<br>Wolf: " + wolf.length + "<br>Sheep: " + sheep.length + "<br>Total: " + (grass.length + wolf.length + sheep.length) + "<br>Time Spent: " + Math.round(simTime / 3600) + " Years";

    clear(); // Makes the little window blank

    fill(255, 255, 0); // Green
    for (let i = 0; i < grass.length; i++) {
        square(grass[i][0]/realLength, grass[i][1]/realLength, width/40);
    }

    fill(255, 128, 0); // Grey
    for (let i = 0; i < wolf.length; i++) {
        square(wolf[i][0]/realLength, wolf[i][1]/realLength, width/40);
    }

    fill(128, 64, 0); // White
    for (let i = 0; i < sheep.length; i++) {
        square(sheep[i][0]/realLength, sheep[i][1]/realLength, width/40);
    }

    let end = Date.now() - t1;
    buckets[parseInt(end/5)]++;
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
    let response = await fetch("http://localhost/remove/" + unqId);
    let data = await response.json();
    unqId = parseInt(Math.random() * 500_000);
    response = await fetch("http://localhost/start/" + unqId);
    data = await response.json();
    simTime = 0;
}

buttonReadout.onclick = async () => {
    const response = await fetch("http://localhost/readout/" + unqId);
    const data = await response.json();
    divReadout.innerHTML = data;
}