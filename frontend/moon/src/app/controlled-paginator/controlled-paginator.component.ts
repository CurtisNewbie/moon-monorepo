import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { PagingController } from 'src/common/paging';

@Component({
  selector: 'app-controlled-paginator',
  templateUrl: './controlled-paginator.component.html',
  styleUrls: ['./controlled-paginator.component.css']
})
export class ControlledPaginatorComponent implements OnInit {

  @ViewChild("paginator", { static: true })
  paginator: MatPaginator;

  @Output("controllerReady")
  controllerEmitter = new EventEmitter<PagingController>();

  pagingController = new PagingController();

  constructor() {

  }

  ngOnInit(): void {
    this.pagingController.control(this.paginator);
    this.controllerEmitter.emit(this.pagingController);
  }

}
