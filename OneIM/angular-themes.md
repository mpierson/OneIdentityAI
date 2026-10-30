# How to create a custom theme for the Identity Manager portal

Pre-requisites: 

- local instance of the Identity Manager API server
- up-to-date copy of the angular source code (branch must match your version of Identity Manager).
- HTML\_qer-app-portal.zip built from latest source deployed to ApiServer
- HotFix for Defect number 651347 may be required

To create the custom theme components, follow the Identity Manager custom theme [README](https://raw.githubusercontent.com/OneIdentity/IdentityManager.Imx/refs/heads/v100/imxweb/custom-theme/readme.md).

Summary of deploy process:  rebuild the CSS with sass 1.90.0, repackage the zip, import it with the Software Loader, then recreate the API server and reload the portal.

**Notes**

- Deployment target for theme files should be "Server\Web\BusinessApiServer"
- Whenever you add or change a theme, users need to reload the page; logging out and back in isn't enough.
- After any import, the API server container must be recreated to pick up the files.
- Start the app server before the API server
- See portal-configuration.md for custom logo information, and configuring a default theme


## Hard-coded styles

Elemental UI and the portal's default stylesheet hardcode colours that a custom Material theme can't reach: the default blue #0a96d1, orange #e36a00 and red #db2534. These hard-coded styles can be superseded with theme-specific styles, e.g. _body.acme-theme_ styles will outrank the default styles and not interfere with other themes.

Example hardcoded styles:

- Top navigation (desktop): the dropdown arrow (.eui-top-navigation-menu-icon) and the active item's text.
- Text and icon colour:
  - Submenu arrows
  - The table toolbar search icon (imx-data-source-toolbar .eui-search-icon)
  - Info alert and banner icons
  - .imx-info, .imx-icon-info, .imx-input-info
  - .imx-checkbox-info frames
  - Hover on .imx-button-icon-hover-info
- Borders:
  - Info alert banners
  - The left edge of info alerts (the right edge in right-to-left layouts)
  - Selected .imx-card-select cards
  - Timeline create and change events
- Backgrounds:
  - .imx-button and hover on .imx-button-hover-info
  - Checked .imx-button-toggle-group-info toggles
  - .imx-chip-info chips
  - The selected side-navigation item
  - div.timeline-customtime
- Both background and border: div.timeline-event-dot.
- Menu items: .mat-mdc-menu-item.imx-primary text and icon colour, set through their Material variables.
- Badges
  - .mat-badge (default)
  - .mat-badge-accent 
  - .mat-badge-warn  
  - .eui-masthead .mat-badge (cart count)
  - .eui-masthead .mat-badge-warn (queue errors)
- Filled buttons
  - --mat-button-filled-container-color
  - --mat-button-filled-label-text-color
- Slide toggles
  - .mat-mdc-slide-toggle / .mat-primary and .mat-accent variants

