var getElementByID = function (id) {
	return document.getElementById(id)
}

var convertBtn = getElementByID("convert_btn")
var inputEle = getElementByID("input")
var outputEle = getElementByID("output")
var warningsEle = getElementByID("warnings_list")

function convert() {
  warningsEle.parentNode.classList.add("hidden")
  outputEle.innerHTML = "converting..."

  while (warningsEle.firstChild) {
    warningsEle.removeChild(warningsEle.firstChild);
  }

  result = curl2httpie.Do(inputEle.value)
  outputEle.innerHTML = result.cmd
  if (result.error !== "") {
    console.log(result.error)
    alert(result.error)
  }
  var warnings = result.warnings
  if (warnings.length > 0) {
    warningsEle.parentNode.classList.remove("hidden")
    var len = warnings.length;
    for (var i = 0; i < len; i++) {
      var li = document.createElement("li")
      var textNode = document.createTextNode(warnings[i])
      li.appendChild(textNode)
      warningsEle.append(li)
    }
  }

  outputEle.select()
}

inputEle.onkeypress = function(evt) {
  evt = evt || window.event;

  if (evt.keyCode == 13) {
    evt.preventDefault()
    convert()
    return false
  }

  return true
}
convertBtn.addEventListener("click", function(event) {
  event.preventDefault()
  convert()
})
