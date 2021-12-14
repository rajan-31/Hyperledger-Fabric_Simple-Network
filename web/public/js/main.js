function doAjax(url, data) {
    return $.ajax({
        method: "POST",
        url: url,
        data: data,
    })
}

$("#submit-create_user").on("click", function(e) {
    e.preventDefault();
    const data = {
        username: $("#username").val(),
        password: $("#password").val(),
        uid: $("#uid").val(),
        new_pass: $("#new_pass").val(),
        new_name: $("#new_name").val(),
    }
    const url = "/api/Create_User"
    
    doAjax(url, data).done(function(res) {
        console.log(res);
        $("#res-create_user").html(JSON.stringify(res));
        $("#create_user").get(0).reset();
    })
    .fail(function(err) {
        console.log(err);
    })
});

$("#submit-getValue").on("click", function(e) {
    e.preventDefault();
    const data = {
        key: $("#key").val()
    }
    const url = "/api/GetValue"
    
    doAjax(url, data).done(function(res) {
        console.log(res);
        $("#res-getValue").html(JSON.stringify(res));
        $("#getValue").get(0).reset();
    })
    .fail(function(err) {
        console.log(err);
    })
})