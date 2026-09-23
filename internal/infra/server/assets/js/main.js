import {DevicesList} from './components.js';

// constants
const pollingInterval = 5;

// gloval variables
const listDevices = new DevicesList(document.getElementById('devices-list-contaier'));
var pollingStatus = null;

// entry point
window.onload = () => {
    console.log("loaded");

    if(!pollingStatus) {
        // start SSE event
        pollingStatus = new EventSource(`/api/devices?interval=${pollingInterval}`);
        // add events
        pollingStatus.onmessage = (event) => {
            console.log('SSE(status) event', 'event=', event);

            // parse data (JSON)
            const data = JSON.parse(event.data);
            if(!data) {
                console.error('json parse error', 'data=', event.data);
                return;
            }

            // update device list
            for(const info of data) {
                listDevices.updateDevice(info);
            }
        }
        pollingStatus.onerror = (event) => {
            console.log('SSE(status) event error', 'event=', event);

            if(pollingStatus) {
                // stop SSE event
                pollingStatus.close();
                pollingStatus = null;
                console.log("stop SSE (polling status)")
            }
        }

        console.log('start SSE (polling status)', 'interval (sec)=', pollingInterval);
    }
}
