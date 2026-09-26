import {DevicesList} from './components.js';
import {PostAdbCommand} from './api.js';

// constants
const pollingInterval = 5;

// gloval variables
const devices = new DevicesList(document.getElementById('devices-list-contaier'));
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
                devices.updateDevice(info);
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

    // add command post event
    document.getElementById('command-post').addEventListener('click', (e) => {
        console.log("command post", "self=",self);

        // clear and hide command error view
        // TODO: implements as compoment
        const view_error = document.getElementById('command-error');
        view_error.textContent = '';
        view_error.classList.remove('show');
        view_error.classList.add('hide');

        // get target serials
        const targets = devices.enabledDevices();
        if(targets.length == 0) {
            view_error.classList.remove('hide');
            view_error.classList.add('show');
            view_error.textContent = 'ERROR: Please select target device(s) more one.';
            return;
        }

        // get args and convert to array
        const view_args = document.getElementById('command-arguments');
        const args = view_args.value.trim().split(/\s+/);
        
        // show result view f all devices
        devices.showAllResult();
        const view_all_clear = document.getElementById('clear-result');
        if(view_all_clear.checked) {
            // clear all result
            devices.clearAllResult();
        }

        // disable post button
        const post_button = e.currentTarget
        post_button.disabled = true;

        // fetch to server
        PostAdbCommand('adb/shell', targets, args, 
            (data, disconnected) => {
                if(disconnected) {
                    // enable post button if disconnected SSE
                    post_button.disabled = false;
                } else {
                    console.log('update command result: ', data);
                    devices.addResult(data.serial, data.result);
                }
            },
            (message) => {
                console.error('command error: ', message);
                view_error.classList.remove('hide');
                view_error.classList.add('show');
                view_error.textContent = 'ERROR: ' + message;
                // enable post button
                post_button.disabled = false;
            }
        )
    });
}
