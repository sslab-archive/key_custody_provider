send_verification_button_event = function (phone) {
    send_verification(phone).then(res => {
        let status = res["status"];
        let jsonRes = res["body"];
        if (status> 300){
            alert("ERR: status -" +status+", body : "+jsonRes);
            return
        }else{
            alert("successfully send verification code to your phone")
        }
        console.log(jsonRes)
    }).catch(err => {
        console.log("ERR!!!");
        console.log(err);
    })
};

check_verification_button_event = function (code) {
    check_verification(code).then(res => {
        console.log(res);
        let status = res["status"];
        let jsonRes = res["body"];
        if (status> 300){
            alert("ERR: status -" +status+", body : "+jsonRes);
            return
        }

        let redirect_url = getParameterByName('redirect_url');
        let r_url = redirect_url + "?purpose=" + jsonRes["purpose"] + "&credential_type=" + jsonRes["credential_type"] +
            "&encrypted_partial_key=" + jsonRes["encrypted_partial_key"] +
            "&encrypted_payload=" + jsonRes["encrypted_payload"] +
            "&partial_key=" + jsonRes["partial_key"] +
            "&partial_key_index=" + jsonRes["partial_key_index"] +
            "&payload=" + jsonRes["payload"] +
            "&provider_id=" + jsonRes["provider_id"] +
            "&public_key=" + jsonRes["public_key"] +
            "&signed_by_private_key=" + jsonRes["signed_by_private_key"];
        console.log(r_url);
        window.location.href = r_url;
    }).catch(err => {
        console.log(err)
    })
};

async function send_verification() {
    let phone = document.getElementById('input_phone').value;
    const sendVerificationResponse = await fetch(
        'http://141.223.121.111:8888/p2/api/authentication/send_code?phone=' + phone,
        {
            method: 'POST'
        }
    ).then(res=>
        res.json().then(data=>({
            status:res.status,
            body:data
        }))
    );
    return await sendVerificationResponse
}


async function check_verification() {
    let phone = document.getElementById('input_phone').value;
    let code = document.getElementById('input_verification_code').value;
    let partial_key = getParameterByName('partial_key');
    let partial_key_index = getParameterByName('partial_key_index');
    let purpose = getParameterByName('purpose');
    let user_public_key = getParameterByName('user_public_key');
    console.log("p:"+purpose);
    const sendVerificationResponse = await fetch(
        'http://141.223.121.111:8888/p2/api/authentication/check?user_public_key='+user_public_key+'&purpose='+purpose+'&phone=' + phone + '&code=' + code + '&partial_key=' + partial_key + '&partial_key_index=' + partial_key_index,
        {
            method: 'POST'
        }
    ).then(res=>
        res.json().then(data=>({
            status:res.status,
            body:data
        }))
    );
    return await sendVerificationResponse
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