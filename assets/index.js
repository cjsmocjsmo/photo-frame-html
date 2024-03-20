$(document).ready(function () {
    $("#fadingImage").hide().fadeIn(5000).delay(25000).fadeOut(5000);
    setTimeout(function () {
        location.reload();
    }, 35000);
    
    // window.location.href = "http://192.168.0.90:8080/update";
});

// $(document).ready(function () {
//     // const x = 1
//     // while (x != 0) {
//         $.ajax({
//             url: "http://192.168.0.90:8080/update",
//             success: function (data) {
//                 $("#fadingImage").attr("src", data).hide().fadeIn(5000).delay(30000).fadeOut(5000);
//             },
//             error: function () {
//                 console.log("Error fetching image");
//             }
//         });
//     }, 35000);