# Mattermost ONLYOFFICE integration plugin
The app which enables the users to edit office documents from Mattermost using ONLYOFFICE Document Server, allows multiple users to collaborate in real-time and to save back those changes to Mattermost

## Features

The app allows to:

* Create and edit text documents, spreadsheets, and presentations.
* Share files with basic permission types - viewing/editing.
* Co-edit documents in real-time.

Supported formats:

* For viewing and editing: DOCX, XLSX, PPTX, CSV, ODT, ODP, ODS.
* For viewing only: DOC, XLS, PPT.

## ONLYOFFICE Docs

You will need an instance of ONLYOFFICE Docs (Document Server) that is resolvable and connectable both from Mattermost and any end clients.

Document Server and Mattermost can be installed either on different computers, or on the same machine. 

If you use one machine, set up a custom port for Document Server as by default both Document Server and ownCloud work on port 80.

You can install free [Community version of ONLYOFFICE Docs](https://helpcenter.onlyoffice.com/installation/docs-community-index.aspx) or scalable [Enterprise Edition with pro features](https://helpcenter.onlyoffice.com/installation/docs-enterprise-index.aspx).

Community Edition vs Enterprise Edition comparison can be found [here](#onlyoffice-docs-editions).

## Installation

1. Clone the [master branch](https://github.com/ONLYOFFICE/onlyoffice-mattermost)
2. Go to the project root
3. Run 
    ```sh
    make dist
    ```
4. Go to <your_mattermost_host>/admin_console/plugins/plugin_management
5. Choose the compiled plugin from your dist folder and press upload

### Plugin settings

- **Document Editing Service address**:
  The URL (and port) of the ONLYOFFICE Docs that provides the editing functionality.

- **Secret key**:
  Is required provided your document server uses JWT security (further information can be found [here] (https://api.onlyoffice.com/editors/signature/))

- **JWT Header**:
  If JWT security is enabled, it is necessary to specify a custom header name since Mattermost's security policy blocks external 'Authorization' Headers. However, this header should be reflected in the ONLYOFFICE Docs signature settings (further information can be found [here] (https://api.onlyoffice.com/editors/signature/))

- **JWT Prefix**:
  Is used to specify the ONLYOFFICE Docs prefix

## ONLYOFFICE Docs editions

ONLYOFFICE offers different versions of its online document editors that can be deployed on your own servers.

* Community Edition (`onlyoffice-documentserver` package)
* Enterprise Edition (`onlyoffice-documentserver-ee` package)

The table below will help you to make the right choice.

| Pricing and licensing | Community Edition | Enterprise Edition |
| ------------- | ------------- | ------------- |
| | [Get it now](https://www.onlyoffice.com/download.aspx)  | [Start Free Trial](https://www.onlyoffice.com/enterprise-edition-free.aspx)  |
| Cost  | FREE  | [Go to the pricing page](https://www.onlyoffice.com/docs-enterprise-prices.aspx)  |
| Simultaneous connections | up to 20 maximum  | As in chosen pricing plan |
| Number of users | up to 20 recommended | As in chosen pricing plan |
| License | GNU AGPL v.3 | Proprietary |
| **Support** | **Community Edition** | **Enterprise Edition** |
| Documentation | [Help Center](https://helpcenter.onlyoffice.com/installation/docs-community-index.aspx) | [Help Center](https://helpcenter.onlyoffice.com/installation/docs-enterprise-index.aspx) |
| Standard support | [GitHub](https://github.com/ONLYOFFICE/DocumentServer/issues) or paid | One year support included |
| Premium support | [Buy Now](https://www.onlyoffice.com/support.aspx) | [Buy Now](https://www.onlyoffice.com/support.aspx) |
| **Services** | **Community Edition** | **Enterprise Edition** |
| Conversion Service                | + | + |
| Document Builder Service          | + | + |
| **Interface** | **Community Edition** | **Enterprise Edition** |
| Tabbed interface                       | + | + |
| Dark theme                             | + | + |
| 150% scaling                           | + | + |
| White Label                            | - | - |
| Integrated test example (node.js)*     | + | + |
| Mobile web editors | - | + |
| Access to pro features via desktop     | - | + |
| **Plugins & Macros** | **Community Edition** | **Enterprise Edition** |
| Plugins                           | + | + |
| Macros                            | + | + |
| **Collaborative capabilities** | **Community Edition** | **Enterprise Edition** |
| Two co-editing modes              | + | + |
| Comments                          | + | + |
| Built-in chat                     | + | + |
| Review and tracking changes       | + | + |
| Display modes of tracking changes | + | + |
| Version history                   | + | + |
| **Document Editor features** | **Community Edition** | **Enterprise Edition** |
| Font and paragraph formatting   | + | + |
| Object insertion                | + | + |
| Adding Content control          | - | + | 
| Editing Content control         | + | + | 
| Layout tools                    | + | + |
| Table of contents               | + | + |
| Navigation panel                | + | + |
| Mail Merge                      | + | + |
| Comparing Documents             | - | +* |
| **Spreadsheet Editor features** | **Community Edition** | **Enterprise Edition** |
| Font and paragraph formatting   | + | + |
| Object insertion                | + | + |
| Functions, formulas, equations  | + | + |
| Table templates                 | + | + |
| Pivot tables                    | + | + |
| Data validation                 | + | + |
| Conditional formatting  for viewing | +** | +** |
| Sheet Views                     | - | + |
| **Presentation Editor features** | **Community Edition** | **Enterprise Edition** |
| Font and paragraph formatting   | + | + |
| Object insertion                | + | + |
| Transitions                     | + | + |
| Presenter mode                  | + | + |
| Notes                           | + | + |
| | [Get it now](https://www.onlyoffice.com/download.aspx)  | [Start Free Trial](https://www.onlyoffice.com/enterprise-edition-free.aspx)  |

\*  It's possible to add documents for comparison from your local drive and from URL.

\** Support for all conditions and gradient. Adding/Editing capabilities are coming soon
