
function postRegister() {
    const username = document.getElementById("username").value
    const email = document.getElementById("email").value
    const password = document.getElementById("password").value

    let reqData = JSON.stringify({
        "username": username,
        "email": email,
        "password": password,
    })

    let xhr = new XMLHttpRequest()
    xhr.onreadystatechange = checkRegisterStatus(xhr)
    xhr.open("post", "/api/register", true)
    xhr.setRequestHeader("Content-Type", "application/json; charset=utf-8")
    xhr.send(reqData)

}

function checkRegisterStatus(xhr) {
    return function() {
        if (xhr.readyState == XMLHttpRequest.DONE &&
            xhr.status == 200) {
            let responseText = xhr.responseText
            let resp = JSON.parse(responseText)
            if (resp["code"] == 0) {
                alert("register success")
                alert(resp["message"])
            } else {
                alert("register failed")
                alert("message: " + resp["message"])
            }
        }
    }
}
