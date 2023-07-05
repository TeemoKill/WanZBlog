function postLogin() {
    const email = document.getElementById("email").value
    const password = document.getElementById("password").value

    let reqData = JSON.stringify({
        "email": email,
        "password": password,
    })

    let xhr = new XMLHttpRequest()
    xhr.onreadystatechange = checkLoginStatus(xhr)
    xhr.open("post", "/api/login", true)
    xhr.setRequestHeader("Content-Type", "application/json; charset=utf-8")
    xhr.send(reqData)

}

function checkLoginStatus(xhr) {
    return function() {
        if (xhr.readyState == XMLHttpRequest.DONE &&
            xhr.status == 200) {
            let responseText = xhr.responseText
            let resp = JSON.parse(responseText)
            if (resp["code"] == 0) {
                alert("login success")
                localStorage.setItem("token", resp["token"])

                // Navigate to user profile page
                window.location.href = "/user/" + resp["user_uuid"]
            } else {
                alert("login failed")
                alert("message: " + resp["message"])
            }
        }
    }
}
