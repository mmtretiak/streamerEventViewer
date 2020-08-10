import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {MatSidenavModule} from "@angular/material/sidenav";
import { RedirectComponent } from './redirect/redirect.component';
import {MatDialogModule} from "@angular/material/dialog";
import {UserService} from "./core/services/user.service";
import {ApiService} from "./core/services/api.service";
import {HTTP_INTERCEPTORS, HttpClient, HttpClientModule} from "@angular/common/http";
import { LoginComponent } from './login/login.component';
import {StreamersService} from "./core/services/streamers.service";
import {ClipService} from "./core/services/clip.service";
import {HomeComponent, SafePipe} from './home/home.component';
import {MatListModule} from "@angular/material/list";
import {MatMenuModule} from "@angular/material/menu";
import {MatIconModule} from "@angular/material/icon";
import {MatButtonModule} from "@angular/material/button";
import {MatToolbarModule} from "@angular/material/toolbar";
import {FlexLayoutModule} from "@angular/flex-layout";
import {HttpTokenInterceptor} from "./core/services/http.token.interceptor";
import {DomSanitizer} from "@angular/platform-browser";
import {AddStreamerDialog} from "./home/add-streamer/add-streamer.dialog";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatInputModule} from "@angular/material/input";
import {FormsModule} from "@angular/forms";
import {NotificationService} from "./core/services/notification.service";
import {MatSnackBar} from "@angular/material/snack-bar";

@NgModule({
  declarations: [
    AppComponent,
    RedirectComponent,
    LoginComponent,
    HomeComponent,
    SafePipe,
    AddStreamerDialog,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatSidenavModule,
    MatDialogModule,
    HttpClientModule,
    MatListModule,
    MatMenuModule,
    MatIconModule,
    MatButtonModule,
    MatToolbarModule,
    FlexLayoutModule,
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
  ],
  providers: [
    UserService,
    ApiService,
    HttpClient,
    StreamersService,
    ClipService,
    { provide: HTTP_INTERCEPTORS, useClass: HttpTokenInterceptor, multi: true },
    MatSnackBar,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
