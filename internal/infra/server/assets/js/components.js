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
        console.debug('Deviceist::findDevice', `serial=`, serial);

        for(const info of this.#list) {
            if(info.serial() == serial) {
                return info;
            }
        }
        return null;
    }
}

export class DeviceCard {
    #view;
    #viewEnabled;
    #viewStatus;

    #serial;
    #status;

    // <li class="device-card" id="device#${this.#serial}">
    //  <input type="checbox" id="device#${this.#serial}-enabled" />
    //  <label for="device#${this.#serial}-enabled">${this.#serial}</label>
    //  <div class="device-model">${info.model}</div>
    //  <div class="device-status ${this.#status}">${this.#status}</div>
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
        viewLabel.for = this.#viewEnabled.id;

        // create div for model
        const viewModel = document.createElement('div');
        viewModel.textContent = info.model;
        viewModel.classList.add('device-model');

        // create div for status
        this.#viewStatus = document.createElement('div');
        this.#viewStatus.textContent = this.#status;
        this.#viewStatus.classList.add('device-status');
        this.#viewStatus.classList.add(this.#status);

        // append children
        this.#view.appendChild(this.#viewEnabled);
        this.#view.appendChild(viewLabel);
        this.#view.appendChild(viewModel);
        this.#view.appendChild(this.#viewStatus);
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
}