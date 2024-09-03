import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

export interface ShareFileQrCodeData {
  title: string;
  msg: string[];
  img: string;
}

@Component({
  selector: 'app-share-file-qrcode-dialog',
  templateUrl: './share-file-qrcode-dialog.component.html',
  styleUrls: ['./share-file-qrcode-dialog.component.css']
})
export class ShareFileQrcodeDialogComponent implements OnInit {

  constructor(
    public dialogRef: MatDialogRef<ShareFileQrcodeDialogComponent, ShareFileQrCodeData>,
    @Inject(MAT_DIALOG_DATA) public data: ShareFileQrCodeData
  ) { }

  ngOnInit(): void {
  }

}
