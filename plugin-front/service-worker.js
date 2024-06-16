console.log("server-worker")
chrome.runtime.onMessage.addListener((message,sender,res) => {
    console.log(message)
    if (message.type === 'set-file') {
        chrome.storage.local.set('file', message.data);
    }
});
chrome.runtime.onMessage.addListener((message,sender,res) => {
    console.log(message)
    if (message.type === 'get-file') {
        chrome.storage.local.get('file',(result)=>{
            res(result)
        });
    }
});