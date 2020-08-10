import { Component, OnInit } from '@angular/core';
import {UserService} from "../core/services/user.service";
import {Router} from "@angular/router";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  constructor(private userService: UserService, private router: Router) {}

  ngOnInit(): void {
    this.userService.login().subscribe((res) => {
      window.location.href = res;
    });
  }
}
