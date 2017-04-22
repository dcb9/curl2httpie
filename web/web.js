var getElementByID = function (id) {
	return document.getElementById(id)
}

var convertBtn = getElementByID("convert_btn")
var inputEle = getElementByID("input")
var outputEle = getElementByID("httpie_output")

convertBtn.addEventListener("click", function(event) {
  outputEle.innerHTML = ""
    outputEle.innerHTML = curl2httpie.Do(inputEle.value)
    event.preventDefault()
})
