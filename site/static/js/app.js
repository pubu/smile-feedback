// JS Goes here - ES6 supported
$('input#submitButton').click( function() {
    $.post( 'https://www.smile-feedback.de/.netlify/functions/service', $('form#emailForm').serialize(), function(data) {
            location.href = data.url;
       },
       'json' // I expect a JSON response
    );
});
