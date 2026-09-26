export async function PostAdbCommand(cmd, targets, args, success, fail) {
    console.debug("post adb command", "cmd=", cmd, "targets=", targets, "args=", args);

    // name URL
    const param = args.join(',') //TODO: join args by ','
    const url = "/api/" + cmd + "?args=" + encodeURI(param);

    try {
        // fetch
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'text/event-stream'
            },
            body: JSON.stringify(targets)
        });

        if (!response.ok) {
            throw new Error(`HTTP Error: ${response.status}`);
        }

        // prepare from stream by UTF-8
        const reader = response.body.getReader();
        const decoder = new TextDecoder('utf-8');
        let buffer = '';

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            // add received chunk to buffer after convert text
            buffer += decoder.decode(value, { stream: true });
            // split SSE message
            const blocks = buffer.split('\n\n');
            // add last un-completed data to buffer
            buffer = blocks.pop();

            for (const block of blocks) {
                if (!block.trim()) continue;

                // split each line
                const lines = block.split('\n');

                // parse event line
                const eventLine = lines.find(l => l.startsWith('event: '));
                const eventType = eventLine ? eventLine.replace(/^event:\s*/, '').trim() : 'message';

                if (eventType === 'close') {
                    console.log('close connection');
                    if(success) { success(null, true); }
                    return; // terminate command execution
                }

                // take data line
                const dataLine= lines.find(l => l.startsWith('data: '));
                if (dataLine) {
                    const rawData = dataLine.replace(/^data:\s*/, '');
                    console.log('received SSE message', "data=", rawData);
                    try {
                        const parsedData = JSON.parse(rawData);
                        
                        if(success) { success(parsedData); }
                    } catch (e) {
                        throw new Error(`json parse error: ${e}`);
                    }
                }
            }
        }
    } catch (error) {
        console.error(error);
        if(fail){ fail(error); }
    }
}
