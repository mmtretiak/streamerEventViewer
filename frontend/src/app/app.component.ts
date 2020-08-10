import {Component, OnInit} from '@angular/core';
import {MatDialog} from "@angular/material/dialog";
import {RedirectComponent} from "./redirect/redirect.component";
import {UserService} from "./core/services/user.service";
import {Router} from "@angular/router";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'twitch-frontend';

  ngOnInit(): void {
    if (!this.userService.isLogged()) {
      localStorage.setItem("token", "token");
      this.router.navigateByUrl("login");
    }
  }

  constructor(private userService: UserService, private router: Router) {}
}
