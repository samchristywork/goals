let currentDateTime = new Date();
let currentDateTimeValue = new Date(currentDateTime.getTime() - currentDateTime.getTimezoneOffset() * 60000).toISOString().slice(0, 16);

document.getElementById("startDateTimePicker").value = currentDateTimeValue;
document.getElementById("endDateTimePicker").value = currentDateTimeValue;

function addGoal() {
  fetch("/add-goal", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: document.getElementById("name").value,
      startAmount: document.getElementById("startAmount").value,
      endAmount: document.getElementById("endAmount").value,
      startTimestamp: document.getElementById("startDateTimePicker").value,
      endTimestamp: document.getElementById("endDateTimePicker").value,
    }),
  }).then(() => {
    window.location.reload();
  });
}

function formatTime(hours) {
  let sign = hours < 0 ? "-" : "";
  hours = Math.abs(hours);
  let h = Math.floor(hours);
  let m = Math.floor((hours - h) * 60);
  let s = Math.floor((hours * 60 - h * 60 - m) * 60);
  return `${sign}${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
}

let globalGoals = null;

function updateGoals() {
  goals = globalGoals;
  if (!goals) {
    return;
  }
  const goalsElement = document.getElementById("goals");

  for (let i=0;i<goals.length;i++) {
    let goal = goals[i];
    let currentTime = new Date().getTime();
    let startTime = new Date(goal.startTimestamp).getTime();
    let endTime = new Date(goal.endTimestamp).getTime();
    let startAmount = parseFloat(goal.startAmount);
    let endAmount = parseFloat(goal.endAmount);
    let currentAmount = parseFloat(goal.currentAmount);

    let timeDiff = endTime - startTime;
    let amountDiff = endAmount - startAmount;
    let percentComplete = (currentTime - startTime) / timeDiff;
    goal.predicted = startAmount + amountDiff * percentComplete;
    goal.ahead = -(goal.predicted - currentAmount) / (amountDiff / timeDiff) / 3600000;
  }

  goals.sort(function(a,b) {
    if (a.ahead < b.ahead) return -1;
    if (a.ahead > b.ahead) return 1;
    return 0;
  });

  html = `
    <table>
      <tr>
        <th>Name</th>
        <th>Start Amt</th>
        <th>End Amt</th>
        <th>Start Time</th>
        <th>End Time</th>
        <th>Current Amt</th>
        <th>Predicted Amt</th>
        <th>Hours Ahead</th>
        <th></th>
      </tr>
  `;

  for(let i=0;i<goals.length;i++) {
    let goal = goals[i];
    html += `
      <tr class="${goal.ahead < 0 ? "behind" : ""}">
        <td onclick="updateField('${goal.name}', 'name', '${goal.name}')">${goal.name}</td>
        <td onclick="updateField('${goal.name}', 'startAmount', '${goal.startAmount}')">${goal.startAmount}</td>
        <td onclick="updateField('${goal.name}', 'endAmount', '${goal.endAmount}')">${goal.endAmount}</td>
        <td onclick="updateField('${goal.name}', 'startTimestamp', '${goal.startTimestamp}')">${goal.startTimestamp.replace("T", " ").replace("Z", "")}</td>
        <td onclick="updateField('${goal.name}', 'endTimestamp', '${goal.endTimestamp}')">${goal.endTimestamp.replace("T", " ").replace("Z", "")}</td>
        <td onclick="updateField('${goal.name}', 'currentAmount', '${goal.currentAmount}')">${goal.currentAmount}</td>
        <td>${goal.predicted.toFixed(2)}</td>
        <td>${formatTime(goal.ahead)}</td>
        <td><button onclick="deleteGoal('${goal.name}')">Delete</button></td>
      </tr>
    `;
  }

  html += "</table>";
  goalsElement.innerHTML = html;
}

function loadGoals() {
  fetch("/goals")
    .then((response) => response.json())
    .then((goals) => {
      globalGoals = goals;
      updateGoals();
    });
}

function updateField(name, column, value) {
  let newValue = prompt(`Enter new ${column}`, value);
  if (newValue === null) {
    return;
  }

  console.log(name, column, value, newValue);

  fetch("/update-field", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: name,
      column: column,
      value: newValue,
    }),
  }).then(() => {
    window.location.reload();
  });
}

function deleteGoal(name) {
  let confirmDelete = confirm("Are you sure you want to delete this goal?");
  if (!confirmDelete) {
    return;
  }

  fetch("/delete-goal", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: name,
    }),
  }).then(() => {
    window.location.reload();
  });
}

loadGoals();
setInterval(updateGoals, 500);
