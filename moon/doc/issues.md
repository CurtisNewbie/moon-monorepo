# Known Issue 

## ngx-lightbox2 rotation (fixed)

ngx-lightbox2 is used in this project. When image rotate, the background may be out-of-sync with the image's rotation, as a result, we can see the white background behind the image. After some investigation, I personally think that this may be caused by the attribute `text-align: center` in class `.lightbox` in `node_modules/ngx-lightbox/lightbox.css` (there may be some conflicts, I am not sure). Removing it somehow fixed the issue :D, so this attribute is overriden in `gallery-image.component.css`.

```css
.lightbox {
  position: absolute;
  left: 0;
  width: 100%;
  z-index: 10000;
  text-align: center;
  line-height: 0;
  font-weight: normal;
  box-sizing: content-box;
  outline: none;
}
```