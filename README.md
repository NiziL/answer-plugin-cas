# CAS Connector for Apache Answer

This is a connector type plugin handling login through CAS.

## Installation

You can follow the [official documentation](https://answer.apache.org/docs/plugins#third-party-plugin) to install this plugin.

```bash
answer build --with github.com/nizil/answer-plugin-cas
```

## Configuration

In the `Admin` panel of Answer, within the `Plugins` menu, enable `CAS Connector`.

Then, from the `CAS Connector` configuration page, you must fill the `Server URL` entry with the URL of your CAS server.
Optionaly, you can provide a `Display Name` to personalize the login button (default to `CAS Connector`).

Finaly, you may want to disable new registration and password login from the `Advanced/Login` panel, hence forcing your users to use the CAS login.

## Usage

Simply login using the new `Connect with <Display Name>` button.
