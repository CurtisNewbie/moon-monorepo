import { Directive, ElementRef, HostListener, Input, OnDestroy, Renderer2 } from '@angular/core';

@Directive({
  selector: '[appImageTooltip]'
})
export class ImageTooltipDirective implements OnDestroy {
  @Input('appImageTooltip') tooltipContent: string;
  @Input() imageUrl: string;
  @Input() imageMaxWidth: string = '200px';

  private tooltipElement: HTMLElement | null = null;

  constructor(
    private elementRef: ElementRef,
    private renderer: Renderer2
  ) {}

  @HostListener('mouseenter')
  onMouseEnter(): void {
    if (!this.imageUrl || this.tooltipElement) {
      return;
    }

    const rect = this.elementRef.nativeElement.getBoundingClientRect();

    const tooltipElement = this.renderer.createElement('div');
    this.renderer.setStyle(tooltipElement, 'position', 'fixed');
    this.renderer.setStyle(tooltipElement, 'background', '#333');
    this.renderer.setStyle(tooltipElement, 'color', '#fff');
    this.renderer.setStyle(tooltipElement, 'padding', '8px');
    this.renderer.setStyle(tooltipElement, 'border-radius', '4px');
    this.renderer.setStyle(tooltipElement, 'box-shadow', '0 2px 8px rgba(0,0,0,0.3)');
    this.renderer.setStyle(tooltipElement, 'z-index', '10000');
    this.renderer.setStyle(tooltipElement, 'max-width', '300px');

    // Position at left edge of row (near thumbnail column), offset by 10px
    this.renderer.setStyle(tooltipElement, 'left', (rect.left + 10) + 'px');

    // Try top first, fall back to bottom if not enough space
    const tooltipHeight = 220; // Approximate height
    const spaceAbove = rect.top;
    const spaceBelow = window.innerHeight - rect.bottom;

    if (spaceAbove >= tooltipHeight + 8) {
      this.renderer.setStyle(tooltipElement, 'top', (rect.top - tooltipHeight - 8) + 'px');
    } else {
      this.renderer.setStyle(tooltipElement, 'top', (rect.bottom + 8) + 'px');
    }

    const img = this.renderer.createElement('img');
    this.renderer.setAttribute(img, 'src', this.imageUrl);
    this.renderer.setStyle(img, 'max-width', this.imageMaxWidth);
    this.renderer.setStyle(img, 'max-height', '200px');
    this.renderer.setStyle(img, 'display', 'block');
    this.renderer.setStyle(img, 'margin-bottom', '4px');

    this.renderer.appendChild(tooltipElement, img);

    if (this.tooltipContent) {
      const text = this.renderer.createElement('div');
      this.renderer.setProperty(text, 'textContent', this.tooltipContent);
      this.renderer.setStyle(text, 'font-size', '12px');
      this.renderer.setStyle(text, 'white-space', 'nowrap');
      this.renderer.setStyle(text, 'overflow', 'hidden');
      this.renderer.setStyle(text, 'text-overflow', 'ellipsis');
      this.renderer.appendChild(tooltipElement, text);
    }

    this.renderer.appendChild(document.body, tooltipElement);
    this.tooltipElement = tooltipElement;
  }

  @HostListener('mouseleave')
  onMouseLeave(): void {
    if (this.tooltipElement) {
      this.renderer.removeChild(document.body, this.tooltipElement);
      this.tooltipElement = null;
    }
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent): void {
    if (this.tooltipElement && !this.elementRef.nativeElement.contains(event.target)) {
      this.renderer.removeChild(document.body, this.tooltipElement);
      this.tooltipElement = null;
    }
  }

  ngOnDestroy(): void {
    if (this.tooltipElement) {
      this.renderer.removeChild(document.body, this.tooltipElement);
      this.tooltipElement = null;
    }
  }
}
