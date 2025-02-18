import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { PdfJsViewerModule } from "ng2-pdfjs-viewer";
import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { MngFilesComponent } from "./manage-files/mng-files.component";
import { MatAutocompleteModule } from "@angular/material/autocomplete";
import { LightboxModule } from "ngx-lightbox";
import {
  APP_BASE_HREF,
  HashLocationStrategy,
  LocationStrategy,
} from "@angular/common";
import { LoginComponent } from "./login/login.component";
import { HttpClientModule, HTTP_INTERCEPTORS } from "@angular/common/http";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { NavComponent } from "./nav/nav.component";
import { ErrorInterceptor } from "./interceptors/error-interceptor";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { MatDatepickerModule } from "@angular/material/datepicker";
import { MatNativeDateModule } from "@angular/material/core";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import {
  MAT_SNACK_BAR_DEFAULT_OPTIONS,
  MatSnackBarModule,
} from "@angular/material/snack-bar";
import { MatTableModule } from "@angular/material/table";
import { MatTooltipModule } from "@angular/material/tooltip";
import { MatPaginatorModule } from "@angular/material/paginator";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatSelectModule } from "@angular/material/select";
import { MatDialogModule } from "@angular/material/dialog";
import { ConfirmDialogComponent } from "./dialog/confirm/confirm-dialog.component";
import { MatMenuModule } from "@angular/material/menu";
import { GrantAccessDialogComponent } from "./grant-access-dialog/grant-access-dialog.component";
import { PdfViewerComponent } from "./pdf-viewer/pdf-viewer.component";
import { ImageViewerComponent } from "./image-viewer/image-viewer.component";
import { MatCheckboxModule } from "@angular/material/checkbox";
import { MatTabsModule } from "@angular/material/tabs";
import { GalleryComponent } from "./gallery/gallery.component";
import { GalleryImageComponent } from "./gallery-image/gallery-image.component";
import { MatCardModule } from "@angular/material/card";
import { FolderComponent } from "./folder/folder.component";
import { MatListModule } from "@angular/material/list";
import { ControlledPaginatorComponent } from "./controlled-paginator/controlled-paginator.component";
import { MediaStreamerComponent } from "./media-streamer/media-streamer.component";
import { TxtViewerComponent } from "./txt-viewer/txt-viewer.component";
import { UserDetailComponent } from "./user-detail/user-detail.component";
import { ManageKeysComponent } from "./manage-keys/manage-keys.component";
import { RegisterComponent } from "./register/register.component";
import { OperateHistoryComponent } from "./operate-history/operate-history.component";
import { ManagerUserComponent } from "./manager-user/manager-user.component";
import { ManageRoleComponent } from "./manage-role/manage-role.component";
import { ManageResourcesComponent } from "./manage-resources/manage-resources.component";
import { ManagePathsComponent } from "./manage-paths/manage-paths.component";
import { MngResDialogComponent } from "./mng-res-dialog/mng-res-dialog.component";
import { MngPathDialogComponent } from "./mng-path-dialog/mng-path-dialog.component";
import { MngRoleDialogComponent } from "./mng-role-dialog/mng-role-dialog.component";
import { ChangePasswordComponent } from "./change-password/change-password.component";
import { AccessLogComponent } from "./access-log/access-log.component";
import { ManageLogsComponent } from "./manage-logs/manage-logs.component";
import { VfolderAddFileComponent } from "./vfolder-add-file/vfolder-add-file.component";
import { HostOnGalleryComponent } from "./host-on-gallery/host-on-gallery.component";
import { DirectoryMoveFileComponent } from "./directory-move-file/directory-move-file.component";
import { ManageBookmarksComponent } from "./manage-bookmarks/manage-bookmarks.component";
import { GalleryAccessComponent } from "./gallery-access/gallery-access.component";
import { ShareFileQrcodeDialogComponent } from "./share-file-qrcode-dialog/share-file-qrcode-dialog.component";
import { MatBadgeModule } from "@angular/material/badge";
import { ListNotificationComponent } from "./list-notification/list-notification.component";
import { BookmarkBlacklistComponent } from "./bookmark-blacklist/bookmark-blacklist.component";
import {
  VerFileHistoryComponent,
  VersionedFileComponent,
} from "./versioned-file/versioned-file.component";
import { TokenInterceptor } from "./interceptors/token-interceptor";
import { CashflowComponent } from "./cashflow/cashflow.component";
import { CashflowStatisticsComponent } from "./cashflow-statistics/cashflow-statistics.component";
import * as PlotlyJS from "plotly.js-dist";
import { PlotlyModule } from "angular-plotly.js";
import { WebpageViewerComponent } from "./webpage-viewer/webpage-viewer.component";
import {
  EditSitePasswordDialogComponent,
  SitePasswordComponent,
  SitePasswordDecryptedDialogComponent,
} from "./site-password/site-password.component";
import { FstoreStorageComponent } from "./fstore-storage/fstore-storage.component";
import { MatTreeModule } from "@angular/material/tree";
import { BrowseHistoryComponent } from "./browse-history/browse-history.component";
import { ScrollingModule } from "@angular/cdk/scrolling";
import { DirTreeNavComponent } from "./dir-tree-nav/dir-tree-nav.component";
import { HIGHLIGHT_OPTIONS, HighlightModule } from "ngx-highlightjs";
import { FileBookmarkDialogComponent } from "./file-bookmark-dialog/file-bookmark-dialog.component";
import { EventPumpDashboardComponent } from './event-pump-dashboard/event-pump-dashboard.component';
import { DroneTaskComponent } from './drone-task/drone-task.component';

PlotlyModule.plotlyjs = PlotlyJS;

@NgModule({
  exports: [],
  declarations: [
    PdfViewerComponent,
    AppComponent,
    MngFilesComponent,
    LoginComponent,
    NavComponent,
    ConfirmDialogComponent,
    GrantAccessDialogComponent,
    ImageViewerComponent,
    GalleryComponent,
    GalleryImageComponent,
    FolderComponent,
    ControlledPaginatorComponent,
    MediaStreamerComponent,
    TxtViewerComponent,
    UserDetailComponent,
    ManageKeysComponent,
    RegisterComponent,
    OperateHistoryComponent,
    ManagerUserComponent,
    ManageRoleComponent,
    ManageResourcesComponent,
    ManagePathsComponent,
    MngResDialogComponent,
    MngPathDialogComponent,
    MngRoleDialogComponent,
    ChangePasswordComponent,
    AccessLogComponent,
    ManageLogsComponent,
    VfolderAddFileComponent,
    HostOnGalleryComponent,
    DirectoryMoveFileComponent,
    ManageBookmarksComponent,
    GalleryAccessComponent,
    ShareFileQrcodeDialogComponent,
    ListNotificationComponent,
    BookmarkBlacklistComponent,
    VersionedFileComponent,
    VerFileHistoryComponent,
    CashflowComponent,
    CashflowStatisticsComponent,
    WebpageViewerComponent,
    SitePasswordComponent,
    SitePasswordDecryptedDialogComponent,
    EditSitePasswordDialogComponent,
    FstoreStorageComponent,
    BrowseHistoryComponent,
    DirTreeNavComponent,
    FileBookmarkDialogComponent,
    EventPumpDashboardComponent,
    DroneTaskComponent,
  ],
  imports: [
    PlotlyModule,
    MatTabsModule,
    MatCheckboxModule,
    MatAutocompleteModule,
    PdfJsViewerModule,
    MatMenuModule,
    BrowserModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    ReactiveFormsModule,
    HttpClientModule,
    MatDatepickerModule,
    MatNativeDateModule,
    FormsModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTableModule,
    MatTooltipModule,
    MatPaginatorModule,
    MatButtonModule,
    MatIconModule,
    MatInputModule,
    MatSelectModule,
    MatDialogModule,
    MatCardModule,
    LightboxModule,
    MatListModule,
    MatBadgeModule,
    MatTreeModule,
    ScrollingModule,
    HighlightModule,
  ],
  entryComponents: [ConfirmDialogComponent, GrantAccessDialogComponent],
  providers: [
    { provide: LocationStrategy, useClass: HashLocationStrategy },
    { provide: APP_BASE_HREF, useValue: "/" },
    { provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true },
    {
      provide: HIGHLIGHT_OPTIONS,
      useValue: {
        fullLibraryLoader: () => import("highlight.js"),
        // lineNumbersLoader: () => import("highlightjs-line-numbers.js"),
      },
    },
    {
      provide: MAT_SNACK_BAR_DEFAULT_OPTIONS,
      useValue: { panelClass: ["snackbarpanel"], duration: 4000 },
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
