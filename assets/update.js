// $(document).ready(function() {
//     function updateImage() {
//         $.ajax({
//             url: 'http://192.168.0.90/images',
//             success: function(data) {
//                 var imageUrl = data.url; // Assuming the response contains the image URL
//                 var img = $('<img>').attr('src', imageUrl).hide();
//                 $('#image-container').empty().append(img);
//                 img.fadeIn(500, function() {
//                     setTimeout(function() {
//                         img.fadeOut(500, function() {
//                             updateImage();
//                         });
//                     }, 30000);
//                 });
//             },
//             error: function() {
//                 // Handle error case
//             }
//         });
//     }

//     updateImage();
// });