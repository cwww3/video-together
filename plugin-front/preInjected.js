console.log('preInject.js');
console.log(window.CW);
(() => {
    if (window.CW) {
        return 
    }

    const originalXMLHttpRequestOpen = XMLHttpRequest.prototype.open;
    XMLHttpRequest.prototype.open = function (...args) {
        try {
            this.addEventListener("load", () => {
                try {
                    const text = this.responseText;
                    const url = this.responseURL;
                    processResponseText(url, text);
                } catch { }
            });
        } catch (e) { console.error(e); }
        return originalXMLHttpRequestOpen.apply(this, args);
    }

    const originalResponseText = Response.prototype.text;
    Response.prototype.text = async function () {
        const text = await originalResponseText.call(this);
        try {
            processResponseText(this.url, text);
        } catch (e) { console.error(e); }
        return text;
    }

    function processResponseText(url, textContent) {
        if (isM3U8(textContent)) {
            const duration = calculateM3U8Duration(textContent);
            console.log({
                'm3u8Url': url,
                'm3u8Content': textContent,
                'duration': String(duration),
            })
            window.postMessage({
                type: "set-file",
                data: {
                    'm3u8Url': url,
                    'm3u8Content': textContent,
                    'duration': String(duration),
                }
            })
        }
    }

    function isM3U8(textContent) {
        return textContent.trim().startsWith('#EXTM3U');
    }

    function calculateM3U8Duration(textContent) {
        let totalDuration = 0;
        const lines = textContent.split('\n');

        for (let i = 0; i < lines.length; i++) {
            if (lines[i].startsWith('#EXTINF:')) {
                let durationLine = lines[i];
                let durationParts = durationLine.split(':');
                if (durationParts.length > 1) {
                    let durationValue = durationParts[1].split(',')[0];
                    let duration = parseFloat(durationValue);
                    if (!isNaN(duration)) {
                        totalDuration += duration;
                    }
                }
            }
        }
        return totalDuration;
    }


})();