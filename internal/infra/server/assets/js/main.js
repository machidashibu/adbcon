import {CommandPalette, DevicesList} from './components.js';
import {PostAdbCommand} from './api.js';

// constants
const pollingInterval = 5;

// gloval variables
const devices = new DevicesList(document.getElementById('devices-list-contaier'));
const command = new CommandPalette(document.getElementById('command-palette'));
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
        console.log("command post");

        // clear and hide command error view
        command.clear();

        // get target serials
        const targets = devices.enabledDevices();
        if(targets.length == 0) {
            command.error('ERROR: Please select target device(s) more one.');
            return;
        }

        // show result view f all devices
        devices.showAllResult();
        // clear all result
        if(command.isClearResult()) {
            devices.clearAllResult();
        }

        // fetch to server
        command.running();
        devices.addAllResult(command.text());
        PostAdbCommand(command.name(), targets, command.args(), 
            (data, disconnected) => {
                if(disconnected) {
                    command.completed();
                } else {
                    console.log('update command result: ', data);
                    devices.addResult(data.serial, data.result);
                }
            },
            (message) => {
                console.error('command error: ', message);
                command.error('ERROR: ' + message);
                command.completed();
            }
        )
    });
}
