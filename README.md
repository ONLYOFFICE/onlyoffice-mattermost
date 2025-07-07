# ONLYOFFICE app for Mattermost
This app enables users to edit office documents from [Mattermost](https://mattermost.com/) using ONLYOFFICE Docs packaged as Document Server - [Community or Enterprise Edition](#onlyoffice-docs-editions).

## Features

The app allows to:

* Open text documents, spreadsheets, presentations, pdf, and diagram.
* Share files with basic permission types - viewing/editing.
* Co-edit documents in real-time.

Supported formats:

**For viewing:**
* **WORD**: DOC, DOCM, DOCX, DOT, DOTM, DOTX, EPUB, FB2, FODT, HTM, HTML, HWP, HWPX, MHT, MHTML, ODT, OTT, PAGES, RTF, STW, SXW, TXT, WPS, WPT, XML
* **CELL**: CSV, ET, ETT, FODS, NUMBERS, ODS, OTS, SXC, XLS, XLSM, XLSX, XLT, XLTM, XLTX
* **SLIDE**: DPS, DPT, FODP, KEY, ODG, ODP, OTP, POT, POTM, POTX, PPS, PPSM, PPSX, PPT, PPTM, PPTX, SXI
* **PDF**: DJVU, DOCXF, OFORM, OXPS, XPS
* **DIAGRAM**: VSDM, VSDX, VSSM, VSSX, VSTM, VSTX

**For editing:**
* **WORD**: DOCM, DOCX, DOTM, DOTX
* **CELL**: XLSB, XLSM, XLSX, XLTM, XLTX
* **SLIDE**: POTM, POTX, PPSM, PPSX, PPTM, PPTX

## Installing ONLYOFFICE Docs

You will need an instance of ONLYOFFICE Docs (Document Server) that is resolvable and connectable both from Mattermost and any end clients. ONLYOFFICE Document Server must also be able to POST to Mattermost directly.

Document Server and Mattermost can be installed either on different computers or on the same machine. In case you choose the latter option, you need to set up a custom port for Document Server. 

You can install free Community version of ONLYOFFICE Docs or scalable Enterprise Edition with pro features.

To install free Community version, use [Docker](https://github.com/onlyoffice/Docker-DocumentServer) (recommended) or follow [these instructions](https://helpcenter.onlyoffice.com/installation/docs-community-install-ubuntu.aspx) for Debian, Ubuntu, or derivatives.

To install Enterprise Edition, follow the instructions [here](https://helpcenter.onlyoffice.com/installation/docs-enterprise-index.aspx).

Community Edition vs Enterprise Edition comparison can be found [here](#onlyoffice-docs-editions).

## Installing ONLYOFFICE app for Mattermost

1. Install Node.js. [Check instructions](https://github.com/nodesource/distributions#installation-instructions)
2. Install Go. [Check instructions](https://go.dev/doc/install)
3. Install make: 

    ```sh
    sudo apt install make
    ```
4. Clone the plugin branch: 

    ```sh
    git clone https://github.com/ONLYOFFICE/onlyoffice-mattermost.git
    ```
5. Go to the project root and start the build:
    ```sh
    cd onlyoffice-mattermost/
    make
    ```

## Plugin settings

![Settings](assets/screen_settings.png)

- **Document Editing Service address**:
  The URL and port of the installed ONLYOFFICE Document Server.

- **Secret key**:
   Starting from version 7.2, JWT is enabled by default and the secret key is generated automatically to restrict the access to ONLYOFFICE Docs and for security reasons and data integrity. Specify your own secret key in the Mattermost plugin configuration. In the ONLYOFFICE Docs [config file](https://api.onlyoffice.com/docs/docs-api/additional-api/signature/), specify the same secret key and enable the validation.

- **JWT Header**:
  If JWT protection is enabled, it is necessary to specify a custom header name since the Mattermost security policy blocks external 'Authorization' Headers. This header should be specified in the ONLYOFFICE Docs signature settings as well (further information can be found [here](https://api.onlyoffice.com/docs/docs-api/additional-api/signature/)).

- **JWT Prefix**:
  Used to specify the ONLYOFFICE Docs prefix.

## Using ONLYOFFICE app for Mattermost

Users are able to open files sent in personal and group Mattermost chats for viewing and co-editing.

When files are sent in the chat message, the following actions are available in the file context menu by clicking the ⋮ symbol: 

- **Open file in ONLYOFFICE** and **Change access rights** — for the author of the message. 
- **Open file in ONLYOFFICE** — for the recipient of the message.

![Settings](assets/screen_actions.png)

When clicking on the _Open file in ONLYOFFICE_ button, the corresponding ONLYOFFICE editor opens in the same window.

![Settings](assets/screen_editor.png)

The author of the message with attached documents is able to change access rights to the file via the context menu using the **Change access rights** option. This action is available both in personal and group chats.

![Settings](assets/screen_share.png)

ONLYOFFICE bot sends notifications about changes in the document specifying the name of the user who made those changes.

![Settings](assets/screen_bot.png)

## ONLYOFFICE Docs editions

ONLYOFFICE offers different versions of its online document editors that can be deployed on your own servers. 

**ONLYOFFICE Docs** packaged as Document Server:

* Community Edition (`onlyoffice-documentserver` package)
* Enterprise Edition (`onlyoffice-documentserver-ee` package)

The table below will help you make the right choice.

| Pricing and licensing | Community Edition | Enterprise Edition |
| ------------- | ------------- | ------------- |
| | [Get it now](https://www.onlyoffice.com/download-docs.aspx?utm_source=github&utm_medium=cpc&utm_campaign=GitHubMattermost#docs-community)  | [Start Free Trial](https://www.onlyoffice.com/download-docs.aspx?utm_source=github&utm_medium=cpc&utm_campaign=GitHubMattermost#docs-enterprise)  |
| Cost  | FREE  | [Go to the pricing page](https://www.onlyoffice.com/docs-enterprise-prices.aspx?utm_source=github&utm_medium=cpc&utm_campaign=GitHubMattermost)  |
| Simultaneous connections | up to 20 maximum  | As in chosen pricing plan |
| Number of users | up to 20 recommended | As in chosen pricing plan |
| License | GNU AGPL v.3 | Proprietary |
| **Support** | **Community Edition** | **Enterprise Edition** |
| Documentation | [Help Center](https://helpcenter.onlyoffice.com/installation/docs-community-index.aspx) | [Help Center](https://helpcenter.onlyoffice.com/installation/docs-enterprise-index.aspx) |
| Standard support | [GitHub](https://github.com/ONLYOFFICE/DocumentServer/issues) or paid | One year support included |
| Premium support | [Contact us](mailto:sales@onlyoffice.com) | [Contact us](mailto:sales@onlyoffice.com) |
| **Services** | **Community Edition** | **Enterprise Edition** |
| Conversion Service                | + | + |
| Document Builder Service          | + | + |
| **Interface** | **Community Edition** | **Enterprise Edition** |
| Tabbed interface                       | + | + |
| Dark theme                             | + | + |
| 125%, 150%, 175%, 200% scaling         | + | + |
| White Label                            | - | - |
| Integrated test example (node.js)      | + | + |
| Mobile web editors                     | - | +* |
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
| Adding Content control          | + | + | 
| Editing Content control         | + | + | 
| Layout tools                    | + | + |
| Table of contents               | + | + |
| Navigation panel                | + | + |
| Mail Merge                      | + | + |
| Comparing Documents             | + | + |
| **Spreadsheet Editor features** | **Community Edition** | **Enterprise Edition** |
| Font and paragraph formatting   | + | + |
| Object insertion                | + | + |
| Functions, formulas, equations  | + | + |
| Table templates                 | + | + |
| Pivot tables                    | + | + |
| Data validation                 | + | + |
| Conditional formatting          | + | + |
| Sparklines                      | + | + |
| Sheet Views                     | + | + |
| **Presentation Editor features** | **Community Edition** | **Enterprise Edition** |
| Font and paragraph formatting   | + | + |
| Object insertion                | + | + |
| Transitions                     | + | + |
| Animations                      | + | + |
| Presenter mode                  | + | + |
| Notes                           | + | + |
| **Form creator features** | **Community Edition** | **Enterprise Edition** |
| Adding form fields              | + | + |
| Form preview                    | + | + |
| Saving as PDF                   | + | + |
| **Working with PDF**      | **Community Edition** | **Enterprise Edition** |
| Text annotations (highlight, underline, cross out) | + | + |
| Comments                        | + | + |
| Freehand drawings               | + | + |
| Form filling                    | + | + |
| | [Get it now](https://www.onlyoffice.com/download-docs.aspx?utm_source=github&utm_medium=cpc&utm_campaign=GitHubMattermost#docs-community)  | [Start Free Trial](https://www.onlyoffice.com/download-docs.aspx?utm_source=github&utm_medium=cpc&utm_campaign=GitHubMattermost#docs-enterprise)  |

\* If supported by DMS.

In case of technical problems, the best way to get help is to submit your issues [here](https://github.com/ONLYOFFICE/onlyoffice-mattermost/issues). Alternatively, you can contact ONLYOFFICE team on [forum.onlyoffice.com](https://forum.onlyoffice.com/).
