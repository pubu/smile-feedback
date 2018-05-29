// JS Goes here - ES6 supported
$('input#submitButton').click( function() {
    $.post( 'https://www.smile-feedback.de/.netlify/functions/service', $('form#emailForm').serialize(), function(data) {
         var img = new Image();
         img.src = "data:image/png;base64," + data.qr;
         $('form#emailForm').append(img);
       },
       'json' // I expect a JSON response
    );
});
