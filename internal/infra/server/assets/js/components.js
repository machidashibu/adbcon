export class DevicesList {
    #view;

    #list = [];

    // <ul class="devices-list">
    //  <li>...</li> <!-- DeviceCard -->
    //  <li>...</li> <!-- DeviceCard -->
    //  ...
    // </ul>
    constructor(container) {
        this.#view = document.createElement('ul');
        this.#view.classList.add('devices-list');

        // append to container
        container.appendChild(this.#view);
    }

    updateDevice(info) {
        console.debug('Deviceist::updateDevice', `info=`, info);

        let device = this.findDevice(info.serial);   // find existing device
        if(device) {
            // update only if existing
            device.status(info.status);
        } else {
            // create new device if not existing
            device = new DeviceCard(info);
            if(!device) {
                console.error('new device create error', 'info=', info);
                return;
            }
            this.#list.push(device);

            // append child
            console.debug('device=', device);
            this.#view.appendChild(device.view());
        }
    }

    findDevice(serial) {
        for(const info of this.#list) {
            if(info.serial() == serial) {
                return info;
            }
        }
        return null;
    }

    enabledDevices() {
        let list = [];
        this.#list.forEach((dev) => {
            if(dev.enabled()) {
                list.push(dev.serial());
            }            
        })
        return list;
    }

    showAllResult() {
        this.#list.forEach((dev) => dev.showResult());
    }

    clearAllResult() {
        this.#list.forEach((dev) => dev.clearResult());
    }

    addResult(serial, text) {
        const device = this.findDevice(serial);   // find existing device
        if(!device) {
            console.error('unknown serial: ', serial);
            return;
        }

        device.addResult(text);
    }

    addAllResult(text) {
        this.#list.forEach((dev) => dev.addResult(text));
    }
}

export class DeviceCard {
    #view;
    #viewEnabled;
    #viewStatus;

    #serial;
    #status;
    #result;

    // <li class="device-card" id="device#${this.#serial}">
    //  <input type="checbox" id="device#${this.#serial}-enabled" />
    //  <label for="device#${this.#serial}-enabled">${this.#serial}</label>
    //  <div class="device-model">${info.model}</div>
    //  <div class="device-status ${this.#status}">${this.#status}</div>
    //  <pre>...</pre> <!-- CommandResult -->
    // </li>
    constructor(info) {
        this.#view = document.createElement('li');
        this.#view.classList.add('device-card');

        this.#serial = info.serial;
        this.#status = info.status;

        // create checkbox for enabled
        this.#viewEnabled = document.createElement('input');
        this.#viewEnabled.type = "checkbox";
        this.#viewEnabled.id = `device#${this.#serial}-enabled`;
        this.#viewEnabled.checked = (this.#status == 'online');  // default checked if status is 'online'
        
        // create label for serial
        const viewLabel = document.createElement('label');
        viewLabel.textContent = this.#serial;
        viewLabel.htmlFor = this.#viewEnabled.id;

        // create div for model
        const viewModel = document.createElement('div');
        viewModel.textContent = info.model;
        viewModel.classList.add('device-model');

        // create div for status
        this.#viewStatus = document.createElement('div');
        this.#viewStatus.textContent = this.#status;
        this.#viewStatus.classList.add('device-status');
        this.#viewStatus.classList.add(this.#status);

        // create view for command result
        this.#result = new CommandResult();

        // append children
        this.#view.appendChild(this.#viewEnabled);
        this.#view.appendChild(viewLabel);
        this.#view.appendChild(viewModel);
        this.#view.appendChild(this.#viewStatus);
        this.#view.appendChild(this.#result.view());
    }

    view() { return this.#view; }
    serial() { return this.#serial; }

    status(status) {
        if(!status) { return this.#status; }    // return status only

        // update status
        console.debug('update status', this.#status, "->", status);
        this.#viewStatus.textContent = status;
        this.#status = status;

        return this.#status;
    }

    enabled(enabled) {
        if(enabled !== Boolean){ return this.#viewEnabled.checked; }    // return enabled only

        // update enabled
        this.#viewEnabled.checked = enabled;

        return this.#viewEnabled.checked;
    }

    addResult(text) { this.#result.add(text); }
    clearResult() { this.#result.clear(); }
    showResult() { this.#result.show(); }
    hideResult() { this.#result.hide(); }
}

export class CommandResult {
    #view;

    // <pre class="console device-console hide"></ul>
    constructor() {
        this.#view = document.createElement('pre');
        this.#view.classList.add("console");
        this.#view.classList.add("device-console");
        this.#view.classList.add("hide");
        this.clear();
    }

    view() { return this.#view; }
    clear() { this.#view.textContent = ''; }
    add(text) { this.#view.textContent += text + '\n'; }
    show() { this.#view.classList.remove('hide'); }
    hide() { this.#view.classList.add('hide'); }
}

export class CommandPalette {
    #viewName;
    #viewArgs;
    #viewPost;
    #viewProgress;
    #viewError;
    #viewClearResult;

    constructor(view) {
        this.#viewName = document.getElementById('command-name');
        this.#viewArgs = document.getElementById('command-arguments');
        this.#viewPost = document.getElementById('command-post');
        this.#viewProgress = document.getElementById('command-progress');
        this.#viewError = document.getElementById('command-error');
        this.#viewClearResult = document.getElementById('clear-result');

        // add selected event to command name
        this.#viewName.addEventListener('change', () => {
            // clear arguments
            this.#viewArgs.value = '';
        })
    }

    clear() {
        this.#viewError.textContent = '';
        this.#viewError.classList.add('hide');
    }
    
    error(message) {
        this.#viewError.classList.remove('hide');
        this.#viewError.textContent = message;
    }

    name() { return this.#viewName.value; }

    args() {
        const value = this.#viewArgs.value.trim();
        if(value) {
            return value.split(/\s+/);
        } else {
            return [];
        }
    }

    text() {
        return '> ' + this.#viewName.options[this.#viewName.selectedIndex].text + ' ' + this.#viewArgs.value;
    }

    running() {
        this.#viewProgress.classList.remove('hide');
        this.lock();
    }

    completed() {
        this.#viewProgress.classList.add('hide');
        this.unlock();
    }

    lock() { this.#viewPost.disabled = true; }
    unlock() { this.#viewPost.disabled = false; }

    isClearResult() { return this.#viewClearResult.checked; }
}