import { Component } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormBuilder, Validators } from '@angular/forms';
import { first } from 'rxjs/operators';

import { AlertService, UserService } from '@app/_services';
import { UserRequest, UserType } from '@app/_models';

@Component({
  templateUrl: 'register.component.html',
})
export class RegisterComponent {
  loading = false;
  submitted = false;
  form = this.formBuilder.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
    type: [1, [Validators.required]],
    imageId: ['test', [Validators.required]],
  });
  UserType = UserType;
  userTypeKeys = Object.keys(this.UserType)
    .filter((k) => !isNaN(Number(k)))
    .map((k) => Number(k));

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
      .register(
        new UserRequest(
          this.form.value.email!,
          this.form.value.password!,
          this.form.value.type!,
          this.form.value.imageId!
        )
      )
      .pipe(first())
      .subscribe({
        next: () => {
          this.alertService.success('Registration successful', {
            keepAfterRouteChange: true,
          });
          this.router.navigate(['../login'], { relativeTo: this.route });
        },
        error: (error) => {
          this.alertService.error(error);
          this.loading = false;
        },
      });
  }
}
