send_verification_button_event = function (email) {
    send_verification(email).then(jsonRes => {
        console.log("success!!");
        console.log(jsonRes)
    }).catch(err => {
        console.log("ERR!!!");
        console.log(err);
    })
};

check_verification_button_event = function (code) {
    check_verification(code).then(jsonRes => {
        console.log(jsonRes);
        let redirect_url = getParameterByName('redirect_url');
        window.location.href = redirect_url + "?credential_type=" + jsonRes["credential_type"] +
            "&encrypted_partial_key=" + jsonRes["encrypted_partial_key"] +
            "&encrypted_payload=" + jsonRes["encrypted_payload"] +
            "&partial_key=" + jsonRes["partial_key"] +
            "&partial_key_index=" + jsonRes["partial_key_index"] +
            "&payload=" + jsonRes["payload"] +
            "&provider_id=" + jsonRes["provider_id"] +
            "&public_key=" + jsonRes["public_key"] +
            "&signed_by_private_key=" + jsonRes["signed_by_private_key"];
    }).catch(err => {
        console.log(err)
    })
};

async function send_verification() {
    let email = document.getElementById('input_email').value;
    const sendVerificationResponse = await fetch(
        'http://127.0.0.1:8888/authentication/send_code?email=' + email
    );
    return await sendVerificationResponse.json()
}


async function check_verification() {
    let email = document.getElementById('input_email').value;
    let code = document.getElementById('input_verification_code').value;
    let partial_key = getParameterByName('partial_key');
    let partial_key_index = getParameterByName('partial_key_index');
    const sendVerificationResponse = await fetch(
        'http://127.0.0.1:8888/api/authentication/check?email=' + email + '&code=' + code + '&partial_key=' + partial_key + '&partial_key_index=' + partial_key_index,
        {
            method: 'POST'
        }
    );
    return await sendVerificationResponse.json()
}

// f56a3370d381c20290f52794f585593d6b7a48ccb4fc0fc5b7ced70353e47507
// 3
function getParameterByName(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, '\\$&');
    var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, ' '));
}