# 위젯 커스터마이징

ALTCHA 위젯은 CSS 변수를 통해 스타일을 변경할 수 있습니다.

## CSS 변수

| 변수 | 기본값 | 설명 |
|---|---|---|
| `--altcha-border-width` | `1px` | 테두리 두께 |
| `--altcha-border-radius` | `3px` | 모서리 둥글기 |
| `--altcha-color-base` | `#ffffff` | 배경색 |
| `--altcha-color-border` | `#a0a0a0` | 테두리 색 |
| `--altcha-color-text` | `currentColor` | 텍스트 색 |
| `--altcha-color-border-focus` | `currentColor` | 포커스 시 테두리 색 |
| `--altcha-color-error-text` | `#f23939` | 에러 텍스트 색 |
| `--altcha-color-footer-bg` | `#f4f4f4` | 푸터 배경색 |
| `--altcha-max-width` | `260px` | 최대 너비 |

## 사용 예시

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

## HTML 속성

| 속성 | 설명 |
|---|---|
| `hidelogo` | 브랜딩 로고 제거 |
| `hidefooter` | 푸터 텍스트 제거 |
| `strings` | JSON 인코딩된 커스텀 라벨/번역 |

```html
<altcha-widget
  challengeurl="/challenge"
  hidelogo
  hidefooter
></altcha-widget>
```

## 참고

- [공식 커스터마이징 문서](https://altcha.org/docs/v2/widget-customization/)
- [Widget v3 (Beta)](https://altcha.org/docs/v2/widget-v3/) — 체크박스, 스위치 토글, 네이티브 체크박스 테마 지원
