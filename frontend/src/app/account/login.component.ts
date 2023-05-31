import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  Validators,
} from '@angular/forms';
import { first } from 'rxjs/operators';

import { AlertService, UserService } from '@app/_services';
import { FormErrors, TokenRequest } from '@app/_models';
import { HttpErrorResponse } from '@angular/common/http';

@Component({ templateUrl: 'login.component.html' })
export class LoginComponent {
  loading = false;
  submitted = false;
  serverError: FormErrors<TokenRequest> | null = null;
  form = this.formBuilder.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', Validators.required],
  });

  constructor(
    private formBuilder: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private userService: UserService,
    private alertService: AlertService
  ) {}

  // convenience getter for easy access to form fields
  get f() {
    return this.form.controls;
  }

  onSubmit() {
    this.submitted = true;

    // reset alerts on submit
    this.alertService.clear();

    // stop here if form is invalid
    if (this.form.invalid) {
      return;
    }

    this.loading = true;
    this.userService
      .login(this.f.email.value!, this.f.password.value!)
      .pipe(first())
      .subscribe({
        next: () => {
          // get return url from query parameters or default to home page
          const returnUrl = this.route.snapshot.queryParams['returnUrl'] || '/';
          this.router.navigateByUrl(returnUrl);
          this.loading = false;
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          this.loading = false;
          this.serverError = error.error as FormErrors<TokenRequest>;
        },
      });
  }
}
