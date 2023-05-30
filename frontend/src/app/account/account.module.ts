import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { CommonModule } from "@angular/common";
import { MatRadioModule } from "@angular/material/radio";

import { AccountRoutingModule } from "./account-routing.module";
import { LayoutComponent } from "./layout.component";
import { LoginComponent } from "./login.component";
import { RegisterComponent } from "./register.component";
import { FormsModule } from "@angular/forms";

@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,
    AccountRoutingModule,
    MatRadioModule,
    FormsModule
  ],
  declarations: [LayoutComponent, LoginComponent, RegisterComponent]
})
export class AccountModule {}
