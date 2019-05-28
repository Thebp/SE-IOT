import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import {ColorPickerModule} from './color-picker/color-picker.module'
import {HttpClientModule} from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations'
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatNativeDateModule} from '@angular/material';
import {MaterialModule} from 'src/app/material/material.module'
import {TreeFlatOverviewExample} from 'src/app/tree-flat-overview-example'
import {platformBrowserDynamic} from '@angular/platform-browser-dynamic';
import {MatSliderModule, MatSlideToggleModule} from '@angular/material'

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    MaterialModule,
    MatNativeDateModule,
    ReactiveFormsModule,
    BrowserModule, 
    ColorPickerModule, 
    HttpClientModule, 
    BrowserAnimationsModule,
    MatSlideToggleModule,
    MatSliderModule
  ],
  declarations: [
    AppComponent,
    TreeFlatOverviewExample,
  ],
  entryComponents: [TreeFlatOverviewExample],
  providers: [],
  bootstrap: [AppComponent, TreeFlatOverviewExample]
})
export class AppModule {

 }
platformBrowserDynamic().bootstrapModule(AppModule);
