export function PostAdbCommand(cmd, targets, args, success, fail) {
    console.debug("post adb command", "cmd=", cmd, "targets=", targets, "args=", args);

    const controller = new AbortController();
    let cancelled = false;

    const cancel = () => {
        cancelled = true;
        controller.abort();
    };

    const run = async () => {
        // name URL
        var url = "/api/" + cmd;
        if(args && args.length != 0) {
            url += "?args=" + encodeURIComponent(args.join(','));
        }

        // make requestInit
        const init = {
            method: 'POST',
            signal: controller.signal
        };
        if(targets.length > 0) {
            init.headers = {
                    'Content-Type': 'application/json',
                    'Accept': 'text/event-stream'
                };
            init.body = JSON.stringify(targets);
        }

        try {
            // fetch
            const response = await fetch(url, init);

            if (!response.ok) {
                const problemDetails = await JSON.parse(rawData);
                throw new Error(`HTTP Error(${problemDetails.status}): ${response.detail}`);
            }

            // prepare from stream by UTF-8
            const reader = response.body.getReader();
            const decoder = new TextDecoder('utf-8');
            let buffer = '';

            while (true) {
                const { done, value } = await reader.read();
                if (done || cancelled) break;

                // add received chunk to buffer after convert text
                buffer += decoder.decode(value, { stream: true });
                // split SSE message
                const blocks = buffer.split('\n\n');
                // add last un-completed data to buffer
                buffer = blocks.pop();

                for (const block of blocks) {
                    if (!block.trim() || cancelled) continue;

                    // parse SSE message
                    const lines = block.split('\n');
                    const eventLine = lines.find(l => l.startsWith('event: '));
                    const eventType = eventLine ? eventLine.replace(/^event:\s*/, '').trim() : 'message';

                    if (eventType === 'close') {
                        console.log('close connection');
                        if (!cancelled && success) success(null, true); // call callback if mot camcelled
                        return; // terminate command execution
                    }

                    // take data line
                    const dataLine= lines.find(l => l.startsWith('data: '));
                    if (dataLine) {
                        const rawData = dataLine.replace(/^data:\s*/, '');
                        console.log('received SSE message', "data=", rawData);
                        try {
                            const parsedData = await JSON.parse(rawData);
                            
                            if(success) { success(parsedData); }
                        } catch (e) {
                            throw new Error(`json parse error: ${e}`);
                        }
                    }
                }
            }
        } catch (error) {
            if (cancelled || error.name === 'AbortError') return;   // not call callback if cancelled

            console.error(error);
            if(fail){ fail(error); }
        }
    };

    run();
    return cancel;
}
