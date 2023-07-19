import { HttpErrorResponse } from '@angular/common/http';
import { Component } from '@angular/core';
import {
  AbstractControl,
  FormControl,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { FormErrors, PostUserDto, UserType } from '@app/_models';
import { AlertService, UserService } from '@app/_services';
import { first } from 'rxjs/operators';

export const passwordMatchingValidatior: ValidatorFn = (
  control: AbstractControl
): ValidationErrors | null => {
  const password = control.get('password');
  const confirmPassword = control.get('confirmPassword');

  return password?.value === confirmPassword?.value
    ? null
    : { notmatched: true };
};

@Component({
  templateUrl: 'register.component.html',
})
export class RegisterComponent {
  loading = false;
  submitted = false;
  serverError: FormErrors<PostUserDto> | null = null;
  profileImage: File | null = null;
  form = new FormGroup(
    {
      email: new FormControl('', [Validators.required, Validators.email]),
      password: new FormControl('', [
        Validators.required,
        Validators.pattern(
          /(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[@$!%*#?&^_-]).{8,}/
        ),
      ]),
      name: new FormControl('', [Validators.required]),
      confirmPassword: new FormControl('', [Validators.required]),
      type: new FormControl(1, [Validators.required]),
      image: new FormControl(null, [Validators.required]),
    },
    { validators: passwordMatchingValidatior }
  );

  UserType = UserType;
  userTypeKeys = Object.keys(this.UserType)
    .filter((k) => !isNaN(Number(k)))
    .map((k) => Number(k));

  constructor(
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
      .register({
        email: this.form.value.email!,
        name: this.form.value.name!,
        password: this.form.value.password!,
        type: this.form.value.type!,
        image: this.profileImage!,
      })
      .pipe(first())
      .subscribe({
        next: () => {
          this.serverError = null;
          this.alertService.success('Registration successful', {
            keepAfterRouteChange: true,
          });
          this.router.navigate(['../login'], { relativeTo: this.route });
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          this.loading = false;
          this.serverError = error.error as FormErrors<PostUserDto>;
        },
      });
  }

  CreateCompareValidator(
    controlOne: AbstractControl,
    controlTwo: AbstractControl
  ) {
    return () => {
      if (controlOne.value !== controlTwo.value)
        return { match_error: 'Value does not match' };
      return null;
    };
  }

  onFileSelect(event: any) {
    if (event.target.files.length > 0) {
      const file = event.target.files[0];
      this.profileImage = file;
    }
  }
}
