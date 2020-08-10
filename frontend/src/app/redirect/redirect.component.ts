import {Component, OnInit} from '@angular/core';
import {MatDialog} from '@angular/material/dialog';
import {ActivatedRoute, Router} from "@angular/router";
import {UserService} from "../core/services/user.service";

/**
 * @title Dialog with header, scrollable content and actions
 */
@Component({
  selector: 'app-login-dialog',
  templateUrl: 'redirect.component.html',
})
export class RedirectComponent implements OnInit {
  constructor(private route: ActivatedRoute, private userService: UserService, private router: Router) {
  }
  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      this.userService.getToken(params.code).subscribe((res) => {
        localStorage.setItem("token", res.token)
        this.router.navigateByUrl('/');
      })
    })
  }
}
