# Widget Customization

The ALTCHA widget supports styling via CSS variables.

## CSS Variables

| Variable | Default | Description |
|---|---|---|
| `--altcha-border-width` | `1px` | Border width |
| `--altcha-border-radius` | `3px` | Border radius |
| `--altcha-color-base` | `#ffffff` | Background color |
| `--altcha-color-border` | `#a0a0a0` | Border color |
| `--altcha-color-text` | `currentColor` | Text color |
| `--altcha-color-border-focus` | `currentColor` | Focus border color |
| `--altcha-color-error-text` | `#f23939` | Error text color |
| `--altcha-color-footer-bg` | `#f4f4f4` | Footer background color |
| `--altcha-max-width` | `260px` | Maximum width |

## Example

```html
<altcha-widget
  challengeurl="http://localhost:3000/challenge"
  style="
    --altcha-border-radius: 8px;
    --altcha-color-border: #d1d5db;
    --altcha-color-border-focus: #3b82f6;
    --altcha-max-width: 320px;
  "
></altcha-widget>
```

## HTML Attributes

| Attribute | Description |
|---|---|
| `hidelogo` | Remove branding logo |
| `hidefooter` | Remove footer text |
| `strings` | JSON-encoded custom labels/translations |

```html
<altcha-widget
  challengeurl="/challenge"
  hidelogo
  hidefooter
></altcha-widget>
```

## References

- [Official Customization Docs](https://altcha.org/docs/v2/widget-customization/)
- [Widget v3 (Beta)](https://altcha.org/docs/v2/widget-v3/) — Checkbox, switch toggle, and native checkbox themes
