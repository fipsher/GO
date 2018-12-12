function init() {
	window.ID = undefined;
	reloadMatrix()
}

function reloadMatrix() {
	let w = $("#width").val();
	let h = $("#height").val();
	$("#matrix").empty();
	let table = $("<table>");
	for (let i = 0; i < h; ++i) {
		let row = $("<tr>");
		for (let j = 0; j < w; ++j) {
			let cell = $("<td>");
			cell.append($("<input>", {
				type: "number",
				class: 'mat',
				row: i,
				col: j
			}));
			row.append(cell);
		}
		if (!i) {
			let col = $("<td>", {
				rowspan: h
			}).html("&times;");
			row.append(col);
		}
		let x = $("<td>");
		x.append($("<input>", {
			disabled: "true",
			class: 'res none',
			row: i,
			value: '?'
		}));
		row.append(x);
		if (!i) {
			let col = $("<td>", {
				rowspan: h
			}).html("=");
			row.append(col);
		}
		let vCell = $("<td>");
		vCell.append($("<input>", {
			type: "number",
			class: 'vec',
			row: i
		}));
		row.append(vCell)
		table.append(row);
	}
	$("#matrix").append(table);
}

function randomMatrix() {
	let w = $("#width").val();
	let max = 10;
	$("#matrix input").each((i, cell) => {
		cell = $(cell);
		if (!cell.hasClass('res')) {
			let row = cell.attr('row');
			let col = cell.attr('col');
			if (row > col) {
				cell.val($(`#matrix input[row=${col}][col=${row}]`).val())
			}
			else {
				let rand = Math.random() * max;
				// need spd for chol method
				// diagonaly dominant is easier to generate
				if (col == row) {
					rand += w * max;
				}
				cell.val(Math.floor(rand));
			}
		}
	});
}

function getMatrix() {
	let arr = $("#matrix input.mat")
		.map((i, cell) => parseFloat(cell.value))
		.toArray();

	let mat = [];
	let n = Math.sqrt(arr.length);
	while(arr.length) {
		mat.push(arr.splice(0, n));
	}
	return mat;
}

function getVector() {
	return $("#matrix input.vec")
		.map((i, cell) => parseFloat(cell.value))
		.toArray();
}

function nAnQ(val) {
	return isNaN(val) ? "?" : val;
}
function setResult(res) {
	if (isNaN(res[0]))
	$("#matrix input.res")
	return $("#matrix input.res")
		.map((i, cell) => {
			cell.value = nAnQ(res[i]);
			cell.title = cell.value;
			$(cell).toggleClass("none", isNaN(res[i]))
		}).toArray();
}

function setProgress(p, text) {
	let pStr = `${p}%`;
	$(".progress .bar").css("width", pStr);
	$(".progress .percent").text(text ? text : pStr);
	$(".progress").toggleClass("stopped", text == "stopped");
}

function getMessage(res) {
	if (res.Status == "cached") {
		// remove timezone mark
		let cached = new Date(res.Timestamp.slice(0, -1));
		let diff =  new Date() - cached;
		let unit = "milliseconds";
		if (diff >= 1000) {
			diff /= 1000;
			diff = Math.floor(diff);
			unit = "second" + (diff == 1 ? "" : "s");
			if (diff >= 60) {
				diff /= 60;
				diff = Math.floor(diff);
				unit = "minute" + (diff == 1 ? "" : "s");
				if (diff >= 60) {
					diff /= 60;
					diff = Math.floor(diff);
					unit = "hour" + (diff == 1 ? "" : "s");
					if (diff >= 24) {
						diff /= 24;
						diff = Math.floor(diff);
						unit = "day" + (diff == 1 ? "" : "s");
					}
				}
			}
		}
		return `cached ${diff} ${unit} ago`;
	}
	return res.Status;
}

function runSolving() {
	if (ID) {
		$.get(`/stop/${ID}`, () => {
			setProgress(0, "stopped");
			ID = undefined;
			$("#run").removeClass("stop");
		})
		return;
	}
	let matrix = JSON.stringify(getMatrix());
	let vector = JSON.stringify(getVector());
	let delay = parseFloat($("#delay").val()) * 1000;
	$.post("/start", {matrix, vector, delay}, id => {
		ID = id;
		$("#run").addClass("stop");
		pullResult();
	});
}

function pullResult() {
	if (!ID) return;
	$.get(`/result/${ID}`, res => {
		setProgress(res.Progress)
		console.log(res);
		if (res.Result) {
			ID = undefined;
			$("#run").removeClass("stop");
			setResult(res.Result);
			setProgress(100, getMessage(res));
		}
		setTimeout(pullResult, 500);
	})
}