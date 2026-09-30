import {CommandPalette, DevicesList} from './components.js';
import {PostAdbCommand} from './api.js';

// constants
const pollingInterval = 5;

// gloval variables
const devices = new DevicesList(document.getElementById('devices-list-contaier'));
const command = new CommandPalette(document.getElementById('command-palette'));
var pollingStatus = null;
let cancelPostAdbCommand = null;

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

        e.preventDefault(); // prevent to reload page

        // clear and hide command error view
        command.clear();

        // show result view f all devices
        // devices.showResult(targets);
        // clear all result
        if(command.isClearResult()) {
            devices.clearResult();
        }

        // get target serials
        let targets = [];
        if(command.isGeneralCommand()) {
            // general command
            command.normal(command.text());
        } else {
            // device command
            targets = devices.enabledDevices();
            if(targets.length == 0) {
                command.error('ERROR: Please select target device(s) more one.');
                return;
            }
            targets.forEach((t) => command.normal(t + command.text()));
        }

        // fetch to server
        command.running();
        cancelPostAdbCommand = PostAdbCommand(command.name(), targets, command.args(), 
            (data, disconnected) => {
                if(disconnected) {
                    cancelPostAdbCommand = null;
                    command.completed();
                } else {
                    console.log('update command result: ', data);
                    if(data.serial != "") {
                        devices.addResult(data.result, data.serial);
                    } else {
                        command.normal(data.result);
                    }
                }
            },
            (message) => {
                cancelPostAdbCommand = null;
                console.error('command error: ', message);
                command.error('ERROR: ' + message);
                command.completed();
            }
        );
    });

    // add command cancel event
    document.getElementById('command-cancel').addEventListener('click', (e) => {
        console.log("command cancel");
        
        e.preventDefault(); // prevent to reload page

        if(cancelPostAdbCommand) {
            cancelPostAdbCommand();
            cancelPostAdbCommand = null;
            command.completed();
        }
    });
}
